package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
    "time"
	"syscall"
)

var (
	url    = flag.String("i", "seashells.io", "URL/IP to use")
	port   = flag.String("p", "1337", "Port to use")
	output = flag.Bool("q", false, "Suppress writing to stdout")
	delay = flag.Int("d", 0, "Delay")
)

func main() {
	flag.Parse()

	fullUrl := *url + ":" + *port
	conn, err := net.Dial("tcp", fullUrl)
	if err != nil {
		log.Fatal(err)
	}

	serverUrl, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Fprint(os.Stderr, serverUrl) // write url to sderr
    
    if *delay > 0 {
        time.Sleep(time.Duration(*delay) * time.Second)
    }

	scan := bufio.NewReader(os.Stdin)
    var both io.Writer
    if *output == true {
        both = conn
    } else {
	    both = io.MultiWriter(os.Stdout, conn) // will write to stdout and the net connection
    }

	done := make(chan error)

	go func() {
		for {
			_, err := syscall.Select(1, &syscall.FdSet{[16]int64{1}}, nil, nil, nil) // check for data on stdin
			done <- err

			_, err = io.Copy(both, scan) // send the scan data to the multiwriter
			done <- err

			_, _, err = scan.ReadLine() // just want to see if we get an EOF
			if err == io.EOF {
				os.Exit(0)
			}
			done <- err
		}
	}()
	for {
		select {
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
