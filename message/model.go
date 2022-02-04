package message

import (
	_ "embed"

	"github.com/zaidanr/go-commander-bot/helper"
)

// TODO: embed
var version string = "v0.0.1"

var (
	UnknownCommand = "Maaf command tidak dikenali."
	HelpMessage    = "*CommanderBot " + version + "*" + "\n\nAvailable Commands:"
)

func Help() *string {
	availCmds := helper.AvailCmds
	msg := HelpMessage
	for _, v := range availCmds {
		msg += "\n==========\nCommand: " + v[0] + "\nDetails: " + v[1] + "\nDescription: " + v[2]
	}
	return &msg
}
