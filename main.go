package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

			fmt.Println(string(buf[:n]))
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

					_, err := conn.Write(body.Bytes())
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
}
