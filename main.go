package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zserge/lorca"
)

func main() {
	var ui lorca.UI
	ui, _ = lorca.New("https://www.baidu.com", "", 800, 600, "--disable-sync", "--disable-translate")
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}
