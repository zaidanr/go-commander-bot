package message

import (
	"fmt"

	"github.com/zaidanr/go-commander-bot/helper"
)

// TODO: embed
var version string = "v0.0.1"

var (
	UnknownCommand = "Maaf command tidak dikenali."
	HelpMessage    = "*CommanderBot " + version + "*" + "\n\nAvailable Commands:"
	Yasss          = "Running *%s* on your command. Please wait.."
	JgnBang        = "Jangan di hek bang, takut 🤕"
)

func Help() *string {
	availCmds := helper.AvailCmds
	msg := HelpMessage
	for _, v := range availCmds {
		msg += "\n==========\nCommand: " + v[0] + "\nDetails: " + v[1] + "\nDescription: " + v[2]
	}
	return &msg
}

func Ack(cmd string) *string {
	msg := fmt.Sprintf(Yasss, cmd)
	return &msg
}
