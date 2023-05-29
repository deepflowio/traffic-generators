package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	"go.uber.org/ratelimit"
)

var (
	fhost    = flag.String("h", "localhost", "Target host")
	fport    = flag.Int("p", 12345, "Target port")
	fsize    = flag.Int("s", 16, "Data size")
	frate    = flag.Int("r", 100000, "Packets per second")
	fthreads = flag.Int("t", 1, "Number of threads")
)

func main() {
	flag.Parse()

	buffer := make([]byte, *fsize)

	// Some random data
	_, err := rand.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < *fthreads; i++ {
		go func(index int) {
			rate := (*frate + *fthreads - 1) / *fthreads
			addr := fmt.Sprintf("%s:%v", *fhost, *fport+index)
			log.Printf("[*] Start flooding %s with data size %d, date rate %d pps.\n", addr, *fsize, rate)

			// connect
			conn, err := net.Dial("udp", addr)
			if err != nil {
				log.Fatal(err)
			}

			// flood
			// Take 10 tokens each time to avoid too high call frequency of the Take() function
			rate_limit := ratelimit.New(rate / 10)
			for {
				rate_limit.Take()
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				conn.Write(buffer)
				//for j := 0; j < 10; j++ {
				//	conn.Write(buffer)
				//}
			}
		}(i)
	}

	// Forever
	<-make(chan bool, 1)
}
