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
	chBackendDie := make(chan struct{})
	chSignal := listenToInterrupt()
	go server.Run(Port)
	go startBrowser(chChromeDie, chBackendDie)
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}
		case <-chChromeDie:
			os.Exit(0)
		}
	}
}

func startBrowser(chChromeDie chan struct{}, chBackendDie chan struct{}) {
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+Port+"/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromeDie <- struct{}{}
	}()
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
