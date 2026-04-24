package player

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Player struct {
	socket string
	conn   *connection
}

func NewPlayer() (p *Player) {
	socket := "/tmp/mpv_rpc"
	var conn *connection
	var cmd *exec.Cmd

	if !mpvExists() {
		log.Fatal("mpv not installed on this machine")
	}

	if supportsKittyProtocol() {
		cmd = exec.Command("mpv", "--vo=kitty", "--idle", "--input-ipc-server="+socket)
	} else {
		cmd = exec.Command("mpv", "--idle", "--input-ipc-server="+socket)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	conn = InitConnection(socket)

	err = conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		socket: socket,
		conn:   conn,
	}
}

func (p *Player) LoadFile(url string) error {
	return p.conn.sendRequest("loadfile", url)
}

func (p *Player) Stop() error {
	return p.conn.sendRequest("stop")
}

func (p *Player) Quit() error {
	return p.conn.sendRequest("quit")
}

func mpvExists() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

func supportsKittyProtocol() bool {
	term := os.Getenv("TERM")
	return strings.Contains(term, "kitty")
}
