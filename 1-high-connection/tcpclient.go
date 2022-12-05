package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	MIN_SERVER_PORT  = 20000
	MAX_SERVER_PORT  = 29999
	NEW_FLOW_PER_SEC = 1000

	REQUEST_PER_FLOW = 18
	REQUEST_INTERVAL = 5 * time.Second
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

func newConnection(dialer *net.Dialer, serverIp string, serverPort int) {
	startTime := time.Now()

	conn, err := dialer.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Printf("Dial server #%d from %s failed ...: %s\n", serverPort, dialer.LocalAddr, err.Error())
		return
	}

	// fmt.Printf("Request %d: local %s, remote %s\n", serverPort, conn.LocalAddr(), conn.RemoteAddr())
	for i := 0; i < REQUEST_PER_FLOW; i++ {
		// read in input from stdin
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//text, _ := reader.ReadString('\n')

		// send to socket
		_, err = fmt.Fprintf(conn, "#"+strconv.Itoa(serverPort)+"-"+strconv.Itoa(i)+"\n")
		if err != nil {
			fmt.Printf("Flow %s->#%d write (%d/%d, %v) failed ...: %s\n",
				dialer.LocalAddr, serverPort, i+1, REQUEST_PER_FLOW, time.Since(startTime), err.Error())
			break
		}

		// listen for reply
		_, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Flow %s->#%d read (%d/%d, %v) failed ...: %s\n",
				dialer.LocalAddr, serverPort, i+1, REQUEST_PER_FLOW, time.Since(startTime), err.Error())
			break
		}
		// fmt.Print("Message from server: " + message)

		time.Sleep(REQUEST_INTERVAL)
	}

	// fmt.Printf("Finished %d: local %s, remote %s\n", serverPort, conn.LocalAddr(), conn.RemoteAddr())
}

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Printf("./tcpclient remote_addr_1,remote_addr_2,... [local_addr_1,local_addr_2...]")
		return
	}
	setLimit()

	serverIps := strings.Split(os.Args[1], ",")
	clientIps := []string{}
	if len(os.Args) == 2 {
		fmt.Printf("Dail %s\n", serverIps)
	} else {
		clientIps = strings.Split(os.Args[2], ",")
		fmt.Printf("Dail %s from %s\n", serverIps, clientIps)
	}

	sleepCounter := 0
	timeStart := time.Now()

	dialer := &net.Dialer{}
	clientIpIndex := 0
	serverIpIndex := 0
	if len(clientIps) > 0 {
		dialer = &net.Dialer{
			LocalAddr: &net.TCPAddr{IP: net.ParseIP(clientIps[clientIpIndex])},
		}
	}

	for serverPort := MIN_SERVER_PORT; ; serverPort++ {
		if serverPort > MAX_SERVER_PORT {
			serverPort = MIN_SERVER_PORT
			if len(clientIps) > 0 {
				clientIpIndex = (clientIpIndex + 1) % len(clientIps)
				dialer = &net.Dialer{
					LocalAddr: &net.TCPAddr{IP: net.ParseIP(clientIps[clientIpIndex])},
				}
			}
			if clientIpIndex == 0 {
				serverIpIndex = (serverIpIndex + 1) % len(serverIps)
			}
		}
		go newConnection(dialer, serverIps[serverIpIndex], serverPort)

		sleepCounter += 1
		if sleepCounter >= NEW_FLOW_PER_SEC {
			timeElapsed := time.Since(timeStart)
			fmt.Printf("Create %d connections, clientIp %s, serverIp %s, serverport %d-%d ~ %d, cost %v\n",
				sleepCounter, clientIps[clientIpIndex], serverIps[serverIpIndex],
				serverPort, sleepCounter, serverPort, timeElapsed)
			if timeElapsed < time.Second {
				time.Sleep(time.Second - timeElapsed)
			}

			timeStart = time.Now()
			sleepCounter = 0
		}
	}
}
