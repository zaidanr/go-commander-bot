package commander

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/sirupsen/logrus"
	"github.com/zaidanr/go-commander-bot/helper"
	"github.com/zaidanr/go-commander-bot/message"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type MyClient struct {
	WAClient *whatsmeow.Client
}

func (cli *MyClient) register() {
	cli.WAClient.AddEventHandler(cli.MessageHandler)
}

func (cli *MyClient) newClient(d *store.Device, l waLog.Logger) {
	cli.WAClient = whatsmeow.NewClient(d, l)
}

func (cli *MyClient) SendMessage(evt interface{}, msg *string) {
	v := evt.(*events.Message)
	resp := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: msg,
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    &v.Info.ID,
				Participant: proto.String(v.Info.MessageSource.Sender.String()),
			},
		},
	}
	cli.WAClient.SendMessage(v.Info.Sender, "", resp)
}

var ClientImpl MyClient
var Log *logrus.Logger

func (cli *MyClient) MessageHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msg := *v.Message.GetExtendedTextMessage().Text

		Log.WithField("From", v.Info.MessageSource.Sender.User).Info(msg)

		msg_arr := strings.Split(strings.ToLower(msg), " ")

		color.Green(msg)

		// Check if commands is predefined in commands.csv
		for i := range helper.AvailCmds {
			if msg_arr[0] == helper.AvailCmds[i][0] {
				var buf bytes.Buffer
				ack := message.Ack(helper.AvailCmds[i][2])
				cli.SendMessage(evt, ack)

				// FIX COMMAND INJECTION
				c := fmt.Sprintf(helper.AvailCmds[i][1], msg_arr[1])
				c_arr := strings.Split(c, " ")

				cmd := exec.Command(c_arr[0], c_arr[1:]...)

				cmd.Stdout = &buf
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				output := buf.String()
				color.Blue(output)
				cli.SendMessage(evt, &output)
				return
			}
		}

		// Additional commands
		switch msg {
		case "/halo":
			cli.SendMessage(evt, proto.String("Halo"))
		case "/test":
			b, _ := exec.Command("whoami").Output()
			out := string(b)
			color.Blue(out)
			cli.SendMessage(evt, &out)
		default:
			cli.SendMessage(evt, message.Help())
		}
		return
	}
}

func init() {
	helper.AvailCmds = helper.ParseCommands()
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:commander.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	ClientImpl.newClient(deviceStore, clientLog)
	ClientImpl.register()

	if ClientImpl.WAClient.Store.ID == nil {
		qrChan, _ := ClientImpl.WAClient.GetQRChannel(context.Background())
		err = ClientImpl.WAClient.Connect()
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}

	} else {
		err = ClientImpl.WAClient.Connect()
		if err != nil {
			panic(err)
		}
	}
}
