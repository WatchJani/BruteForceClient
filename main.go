package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer conn.Close()
	done := make(chan struct{})

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println(err)

				if err == io.EOF {
					break //kill the go routine
				}
			}

			fmt.Println("[server]: ", string(buf[:n]))
		}

		fmt.Println("listener close")
		close(done)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		for {
			select {
			case <-done:
				fmt.Println("Connection is closed")
				return
			default:
				if scanner.Scan() {
					text := scanner.Text()
					if text == "quit" {
						fmt.Println("Client close connection")
						return
					}

					body := bytes.NewBufferString(text)
					req, err := Parser(body.String())
					if err != nil {
						log.Println(err)
						continue
					}

					fmt.Println(req)

					_, err = conn.Write([]byte(req))
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
}

//start hash [single/thread/multiple]

func Parser(input string) (string, error) {
	cmd := strings.Split(input, " ")

	switch cmd[0] {
	case "start":
		return Start(cmd[1:])
	case "end":
		return End(cmd[1:])
	default:
		return "", fmt.Errorf("wrong command")
	}
}

func End(args []string) (string, error) {
	return "", nil
}

type State struct {
	Hash string `json:"hash"` //Delete, put in Master
	Mod  string `json:"mod"`  //Delete
}

func Start(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("wrong parameters")
	}

	fmt.Println(args)

	state := State{
		Hash: args[0],
		Mod:  args[1],
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("cmd: start\nbody: %s", data), nil
}
