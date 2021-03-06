package main

import (
	"io"
	"log"
	"net"
)

const (
	comQuit byte = iota + 1
	comInitDB
	comQuery
	comFieldList
	comCreateDB
	comDropDB
	comRefresh
	comShutdown
	comStatistics
	comProcessInfo
	comConnect
	comProcessKill
	comDebug
	comPing
	comTime
	comDelayedInsert
	comChangeUser
	comBinlogDump
	comTableDump
	comConnectOut
	comRegiserSlave
	comStmtPrepare
	comStmtExecute
	comStmtSendLongData
	comStmtClose
	comStmtReset
	comSetOption
	comStmtFetch
)

func main() {
	ln, err := net.Listen("tcp", ":3316")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go proxify(conn)
	}
}

func proxify(conn net.Conn) {
	server, err := net.Dial("tcp", "localhost:3306")
	if err != nil {
		log.Println("Could not dial server")
		log.Println(err)
		conn.Close()
		return
	}
	go forward(server, conn)
	forwardWithLog(conn, server)
	server.Close()
	conn.Close()
}

func forward(src, sink net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return
		}

		_, err = sink.Write(buffer[0:n])
		if err != nil && err != io.EOF {
			return
		}
	}
}

func forwardWithLog(src, sink net.Conn) {
	buffer := make([]byte, 16777219)
	for {
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return
		}

		if n >= 5 {
			switch buffer[4] {
			case comQuery:
				log.Printf("Query: %s\n", string(buffer[5:n]))
			case comStmtPrepare:
				log.Printf("Prepare Query: %s\n", string(buffer[5:n]))
			}

			switch buffer[11] {
			case comQuery:
				log.Printf("Query: %s\n", string(buffer[12:n]))
			case comStmtPrepare:
				log.Printf("Prepare Query: %s\n", string(buffer[12:n]))
			}
		}

		_, err = sink.Write(buffer[0:n])
		if err != nil && err != io.EOF {
			return
		}
	}

}
