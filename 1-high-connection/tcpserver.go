package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

// import "strings"

const (
	MIN_SERVER_PORT = 10000
	MAX_SERVER_PORT = 29999
)

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	fmt.Printf("set cur limit: %d\n", rLimit.Cur)
}

func main() {
	setLimit()
	liveConnection := int64(0)
	totalConnection := uint64(0)
	prevConnection := uint64(0)
	totalPacket := uint64(0)
	prevPacket := uint64(0)

	// listen on all interfaces
	for port := MIN_SERVER_PORT; port <= MAX_SERVER_PORT; port++ {
		ln, listenErr := net.Listen("tcp", ":"+strconv.Itoa(port))
		if listenErr != nil {
			fmt.Printf("Launching server %d (total = %d) failed ...: %s\n", port, atomic.LoadUint64(&totalConnection), listenErr.Error())
			continue
		}

		// accept connection on port
		go func(server net.Listener, serverPort int) {
			for {
				conn, acceptErr := server.Accept()
				if acceptErr != nil {
					fmt.Printf("Server #%d (total = %d) accept failed ...: %s\n", serverPort, atomic.LoadUint64(&totalConnection), acceptErr.Error())
					continue
				}
				atomic.AddInt64(&liveConnection, 1)
				atomic.AddUint64(&totalConnection, 1)
				go func(serverPort int) {
					for i := 0; ; i++ {
						// will listen for message to process ending in newline (\n)
						message, err := bufio.NewReader(conn).ReadString('\n')
						if err != nil {
							//fmt.Printf("Server #%d read %d failed ...: %s\n", serverPort, i, err.Error())
							break
						}
						// output message received
						//fmt.Print("Message Received:", string(message))

						// sample process for string received
						//newmessage := strings.ToUpper(message)

						// send new string back to client
						_, err = conn.Write([]byte(message + "\n"))
						if err != nil {
							//fmt.Printf("Server #%d write %d failed ...: %s\n", serverPort, i, err.Error())
							break
						}
						atomic.AddUint64(&totalPacket, 1)
					}
					atomic.AddInt64(&liveConnection, -1)
					// fmt.Printf("Finished: local %s, remote %s, live %d, total %d\n", conn.LocalAddr(), conn.RemoteAddr(), atomic.LoadInt64(&liveConnection), atomic.LoadUint64(&totalConnection))
				}(serverPort)
			}
			fmt.Printf("Server #%d exit ...: %s\n", serverPort)
		}(ln, port)
	}

	ticker := time.NewTicker(1 * time.Second)
	i := 0
	for range ticker.C {
		i += 1
		currConnection := atomic.LoadUint64(&totalConnection)
		currPacket := atomic.LoadUint64(&totalPacket)
		newConnection := currConnection - prevConnection
		fmt.Printf("%s %d liveConn %d, newConn %d, totalConn %d, pps %d\n",
			time.Now().Format(time.RFC3339), i,
			atomic.LoadInt64(&liveConnection), newConnection, currConnection,
			currPacket-prevPacket)
		prevConnection = currConnection
		prevPacket = currPacket
	}
}
