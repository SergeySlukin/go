package main

import (
	"flag"
	"os"
	"fmt"
	"time"
	"net"
	"strings"
	"encoding/hex"
)

var (
	host *string = flag.String("host", "", "target host or address")
	port *string = flag.String("port", "0", "target port")
	listenPort *string = flag.String("listen_port", "0", "listen port")
)

func die(format string, v ...interface{})  {
	os.Stderr.WriteString(fmt.Sprintf(format + "\n", v...))
	os.Exit(1)
}

func connectionLogger(data chan []byte, connNumber int, localInfo, remoteInfo string)  {
	logName := fmt.Sprintf("log-%s-%04d-%s-%s.log", formatTime(time.Now()), connNumber, localInfo, remoteInfo)
	loggerLoop(data, logName)
}

func binaryLogger(data chan []byte, connNumber int, peer string) {
	logName := fmt.Sprintf("log-binary-%s-%04d-%s.log", formatTime(time.Now()), connNumber, peer)
	loggerLoop(data, logName)
}

func loggerLoop(data chan []byte, logName string)  {
	f, err := os.Create(logName)
	if err != nil {
		die("Unable to create file %s, %v\n", logName, err)
	}
	defer f.Close()
	for {
		b := <-data
		if len(b) == 0 {
			break
		}
		f.Write(b)
		f.Sync()
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006.01.02-15.04.05")
}

func printableAddr(a net.Addr) string  {
	return strings.Replace(a.String(), ":", "-", -1)
}

type Channel struct {
	from, to net.Conn
	logger, binaryLogger chan []byte
	ack chan bool
}

func passThrough(c *Channel)  {
	fromPeer := printableAddr(c.from.LocalAddr())
	toPeer := printableAddr(c.to.RemoteAddr())

	b := make([]byte, 10240)
	offset := 0
	packetN := 0
	for {
		n, err := c.from.Read(b)
		if err != nil {
			c.logger <- []byte(fmt.Sprintf("Disconnected from %s \n", fromPeer))
			break
		}
		if n > 0 {
			c.logger <- []byte(fmt.Sprintf("Received ($%d, %08X) %d bytes from %s\n", packetN, offset, n, fromPeer))

			c.logger <- []byte(hex.Dump(b[:n]))
			c.binaryLogger <- b[:n]
			c.to.Write(b[:n])
			c.logger <- []byte(fmt.Sprintf("Sent (#%d) to %s\n", packetN, toPeer))
			offset += n
			packetN += 1
		}
	}
	c.from.Close()
	c.to.Close()
	c.ack <- true
}

func processConnection(local net.Conn, connNumber int, target string) {
	remote, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Printf("Unable to connect to %s, %v\n", target, err)
	}

	localInfo := printableAddr(remote.LocalAddr())
	remoteInfo := printableAddr(remote.RemoteAddr())

	started := time.Now()

	logger := make(chan []byte)
	fromLogger := make(chan []byte)
	toLogger := make(chan []byte)

	ack := make(chan bool)
	go connectionLogger(logger, connNumber, localInfo, remoteInfo)
	go binaryLogger(fromLogger, connNumber, localInfo)
	go binaryLogger(toLogger, connNumber, remoteInfo)

	logger <- []byte(fmt.Sprintf("Connected to %s at %s\n", target, formatTime(started)))

	go passThrough(&Channel{remote, local, logger, toLogger, ack})
	go passThrough(&Channel{local, remote, logger, fromLogger, ack})

	<-ack
	<-ack

	finished := time.Since(started)
	logger <- []byte(fmt.Sprintf("Finished at %s, duration %s\n",
		formatTime(started), finished.Seconds()))
	logger <- []byte{}
	fromLogger <- []byte{}
	toLogger <- []byte{}
}

func main()  {
	flag.Parse()

	if flag.NFlag() != 3{
		fmt.Printf("usage: gotcpspy -host target_host -port target_port -listen_post=local_port\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	target := net.JoinHostPort(*host, *port)
	fmt.Printf("Start listening on port %s and forwarding data to %s\n",
		*listenPort, target)
	ln, err := net.Listen("tcp", ":"+*listenPort)
	if err != nil {
		fmt.Printf("Unable to start listener, %v\n", err)
		os.Exit(1)
	}
	connN := 1
	for {
		if conn, err := ln.Accept(); err == nil {
			go processConnection(conn, connN, target)
			connN++
		} else {
			fmt.Printf("Accept failed, %v\n", err)
		}
	}
}
