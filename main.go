package main

import (
	"os"
	"os/exec"
	"os/signal"

	"example.com/m/server"
)

var Port = "27149"

func main() {
	chChromeDie := make(chan struct{})
	go server.Run(Port)
	go startBrowser(chChromeDie)
	chSignal := listenToInterrupt()
	select {
	case <-chSignal:
	case <-chChromeDie:
		os.Exit(0)
	}
}

func startBrowser(chChromeDie chan struct{}) {
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+Port+"/static/index.html")
	cmd.Start()
	cmd.Wait()
	chChromeDie <- struct{}{}
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
