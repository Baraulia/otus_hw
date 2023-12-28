package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	address, timeout := parseFlags()
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	closeChan := make(chan struct{}, 1)
	errorChanel := make(chan error, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	wg.Add(3)
	go handleReceiving(&wg, client, errorChanel, closeChan)
	go handleSending(&wg, client, errorChanel, closeChan)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-sigChan:
				log.Println("getting sigkill signal")
				client.Close()
				close(closeChan)
				return
			case <-errorChanel:
				client.Close()
				close(closeChan)
				return
			}
		}
	}()

	wg.Wait()
	close(sigChan)
	close(errorChanel)
}

func parseFlags() (string, time.Duration) {
	timeout := pflag.Duration("timeout", 10*time.Second, "timeout")
	pflag.Parse()
	host := pflag.Arg(0)
	port := pflag.Arg(1)

	if host == "" || port == "" {
		fmt.Println("you need to specify host and port ")
		os.Exit(1)
	}

	address := net.JoinHostPort(host, port)

	return address, *timeout
}

func handleSending(wg *sync.WaitGroup, client TelnetClient, errorChan chan error, closeChan chan struct{}) {
	defer wg.Done()
	go func() {
		err := client.Send()
		if err != nil {
			errorChan <- err
			if errors.Is(err, io.EOF) {
				log.Println("...EOF")
				return
			}
			return
		}
	}()

	for range closeChan {
		log.Println("... finishing goroutine with send")
		return
	}
}

func handleReceiving(wg *sync.WaitGroup, client TelnetClient, errorChan chan error, closeChan chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-closeChan:
			log.Println("... finishing goroutine with receive")
			return
		default:
			err := client.Receive()
			if err != nil {
				errorChan <- err
				if errors.Is(err, io.EOF) {
					log.Println("...Connection was closed by peer")
					return
				}
				return
			}
		}
	}
}
