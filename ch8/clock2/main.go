package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

var port = flag.Uint64("port", 8000, "Listen port")
var timezone = flag.String("tz", "Local", "Time zone")

func main() {
	// get command lines flags
	flag.Parse()
	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.FormatUint(*port, 10))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Listening ", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := io.WriteString(c, time.Now().In(loc).Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
