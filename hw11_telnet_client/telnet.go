package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	address string
	timeout time.Duration
	input   io.ReadCloser
	output  io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		input:   in,
		output:  out,
	}
}

func (t *Telnet) Connect() error {
	log.Printf("...Connected to %s", t.address)
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("error while connecting: %w", err)
	}

	t.conn = conn

	return nil
}

func (t *Telnet) Close() error {
	log.Println("...Closing connection")
	err := t.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Send() error {
	reader := bufio.NewReader(t.input)
	data := make([]byte, 1024)
	n, err := reader.Read(data)
	if err != nil {
		return err
	}

	_, err = t.conn.Write(data[:n])
	if err != nil {
		return fmt.Errorf("error while writening: %w", err)
	}

	return nil
}

func (t *Telnet) Receive() error {
	buf := make([]byte, 1024)
	n, err := t.conn.Read(buf)
	if err != nil {
		return err
	}

	_, err = t.output.Write(buf[:n])
	if err != nil {
		return err
	}

	return nil
}
