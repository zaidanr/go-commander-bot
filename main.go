package main

import (
	_ "embed"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zaidanr/go-commander-bot/commander"
	_ "github.com/zaidanr/go-commander-bot/commander"
	"github.com/zaidanr/go-commander-bot/helper"
)

func main() {
	helper.AvailCmds = helper.ParseCommands()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	commander.ClientImpl.WAClient.Disconnect()
}
