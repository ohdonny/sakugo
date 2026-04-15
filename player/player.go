package player

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Connection struct {
	client     net.Conn
	socketPath string
	error      string
	events     chan Event
}

type CommandRequest struct {
	Command []any `json:"command"`
}

type CommandResponse struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}

type Event struct {
	Name   string `json:"event"`
	Reason string `json:"reason"`
}

func InitConnection(socketPath string) *Connection {
	return &Connection{
		socketPath: socketPath,
		events:     make(chan Event, 1),
	}
}

func (c *Connection) Open() error {
	if c.client != nil {
		return fmt.Errorf("client already open")
	}
	client, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return fmt.Errorf("could not find path to MPV")
	}
	c.client = client

	go c.listen()
	return nil
}

func (c *Connection) Close() error {
	err := c.client.Close()
	c.client = nil
	return err
}

func (c *Connection) LoadFile(url string) error {
	return c.sendRequest("loadfile", url)
}

func (c *Connection) Stop() error {
	return c.sendRequest("stop")
}

func (c *Connection) Quit() error {
	return c.sendRequest("quit")
}

func (c *Connection) sendRequest(commands ...any) error {
	message := CommandRequest{
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

func (c *Connection) readEvent(data []byte) {
	var event Event

	err := json.Unmarshal(data, &event)

	if err != nil {
		return
	}

	if event.Name == "end-file" {
		c.events <- event
	}
}

func (c *Connection) readResponse(data []byte) {
	var response CommandResponse
	err := json.Unmarshal(data, &response)

	if err != nil {
		return
	}

	if response.Error == "" || response.Error == "success" {
		return
	}
}

func (c *Connection) listen() {
	scanner := bufio.NewScanner(c.client)
	for scanner.Scan() {
		data := scanner.Bytes()
		c.readEvent(data)
		c.readResponse(data)
	}
}
