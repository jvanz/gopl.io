package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

const sep = "="

var mu sync.Mutex
var times = make(map[string]string)

func main() {
	servers := parseArguments()
	for city, server := range servers {
		go getTime(city, server)
	}
	printTime()
}

func parseArguments() map[string]string {
	if len(os.Args) < 2 {
		log.Fatal("At least one remote server required")
	}
	var servers = make(map[string]string)
	for _, arg := range os.Args[1:] {
		splits := strings.Split(arg, sep)
		servers[splits[0]] = splits[1]
	}
	return servers
}

func printTime() {
	const format = "%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// print time table header
	fmt.Fprintf(tw, format, "City", "Time")
	fmt.Fprintf(tw, format, "-----", "-----")
	for {
		mu.Lock()
		if len(times) > 1 {
			for city, time := range times {
				fmt.Fprintf(tw, format, city, time)
			}
		}
		mu.Unlock()
		time.Sleep(1000 * time.Millisecond)
		tw.Flush()

	}
}

func getTime(city, server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		mu.Lock()
		times[city] = scanner.Text()
		mu.Unlock()
		time.Sleep(1000 * time.Millisecond)
	}
}
