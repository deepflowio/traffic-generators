package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/test", HandlePost)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// ticker := time.NewTicker(time.Second)
	// go func() {
	// 	for t := range ticker.C {
	// 		io.WriteString(w, "Chunk")
	// 		fmt.Println("Tick at", t)
	// 	}
	// }()
	// time.Sleep(time.Second * 30)
	// ticker.Stop()
	for i := 0; i < 30; i++ {
		io.WriteString(w, "Chunk")
		fmt.Println("Tick at", i)
		time.Sleep(time.Second)
	}
	fmt.Println("Finished: should return Content-Length: 0 here")
	w.Header().Set("Content-Length", "0")
}
