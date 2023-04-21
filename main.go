package main

import (
	"os"
	"os/exec"
	"os/signal"

	"example.com/m/server"
)

var Port = "27149"

func main() {
	go server.Run(Port)
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+Port+"/static")
	cmd.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
