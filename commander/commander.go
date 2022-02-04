package commander

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
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

func (cli *MyClient) MessageHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msg := *v.Message.GetExtendedTextMessage().Text

		color.Green(msg)

		for i := range helper.AvailCmds {
			if msg == helper.AvailCmds[i][0] {
				cli.SendMessage(evt, &msg)
				return
			}
		}

		if msg == "/Halo" {
			cli.SendMessage(evt, proto.String("Halo"))
			return
		}

		cli.SendMessage(evt, message.Help())
	}
}

var ClientImpl MyClient

func init() {
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
