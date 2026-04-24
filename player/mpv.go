package player

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type connection struct {
	client     net.Conn
	socketPath string
	error      string
	events     chan event
}

type commandRequest struct {
	Command []any `json:"command"`
}

type commandResponse struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}

type event struct {
	Name   string `json:"event"`
	Reason string `json:"reason"`
}

func InitConnection(socketPath string) *connection {
	return &connection{
		socketPath: socketPath,
		events:     make(chan event, 1),
	}
}

func (c *connection) Open() error {
	if c.client != nil {
		return fmt.Errorf("client already open")
	}

	maxRetries := 5
	var client net.Conn
	var err error

	for i := range maxRetries {
		client, err = net.Dial("unix", c.socketPath)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * 50 * time.Millisecond)
	}

	if err != nil {
		return fmt.Errorf("could not find path to MPV")
	}

	c.client = client

	go c.listen()
	return nil
}

func (c *connection) Close() error {
	err := c.client.Close()
	c.client = nil
	return err
}

func (c *connection) sendRequest(commands ...any) error {
	message := commandRequest{
		Command: commands,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = c.client.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) readEvent(data []byte) {
	var event event

	err := json.Unmarshal(data, &event)

	if err != nil {
		return
	}

	if event.Name == "end-file" {
		c.events <- event
	}
}

func (c *connection) readResponse(data []byte) {
	var response commandResponse
	err := json.Unmarshal(data, &response)

	if err != nil {
		return
	}

	if response.Error == "" || response.Error == "success" {
		return
	}
}

func (c *connection) listen() {
	scanner := bufio.NewScanner(c.client)
	for scanner.Scan() {
		data := scanner.Bytes()
		c.readEvent(data)
		c.readResponse(data)
	}
}
