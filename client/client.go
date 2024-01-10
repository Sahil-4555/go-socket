package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var availableGroups = map[int]string{
	1: "ABC (GROUP-1)",
	2: "AB (GROUP-2)",
}

type Data struct {
	Username string
	Message  string
}

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// ANSI escape codes for text colors and cursor movement
const (
	Blue      = "\033[34m"
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	CursorUp  = "\033[1A"
	ClearLine = "\033[2K"
)

func listenForMessages(conn net.Conn, myColor string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Disconnected from server.")
			conn.Close()
			os.Exit(1)
		}
		fmt.Print(CursorUp + ClearLine)
		fmt.Println(myColor + Bold + message + Reset)
		fmt.Print("Text to send: ")
	}
}

func main() {
	fmt.Println("Connecting to", connType, "server", connHost+":"+connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n")
	for i, v := range availableGroups {
		fmt.Printf("%d. %s\n", i, v)
	}

	fmt.Printf("Enter Into the Group: ")
	var group int
	fmt.Scan(&group)
	conn.Write([]byte(availableGroups[group] + "\n"))
	myColor, _ := bufio.NewReader(conn).ReadString('\n')
	myColor = myColor[:len(myColor)-1]

	go listenForMessages(conn, myColor)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Text to send: ")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
		}
		tmp := Data{Message: input, Username: username}

		encoder := gob.NewEncoder(conn)
		err = encoder.Encode(&tmp)
		if err != nil {
			fmt.Println("Error Decoding: ", err.Error())
			break
		}
		fmt.Print(CursorUp + ClearLine)
		fmt.Println("You: " + input[:len(input)-1]) // Remove newline from input
	}
}
