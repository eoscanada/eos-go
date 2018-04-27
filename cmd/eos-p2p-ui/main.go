package main

import (
	"bufio"
	"os"

	"fmt"

	"log"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/p2p"
	"github.com/marcusolsson/tui-go"
)

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "john", message: "hi, what's up?", time: "14:41"},
	{username: "jane", message: "not much", time: "14:43"},
}

var blockLog *tui.Box
var transactionLog *tui.Box
var ui tui.UI

func main() {

	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			//do nothing
		}
	}()

	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("general"),
		tui.NewLabel("random"),
		tui.NewLabel(""),
		tui.NewLabel("DIRECT MESSAGES"),
		tui.NewLabel("slackbot"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	blockLog = tui.NewVBox()

	blockLogScroll := tui.NewScrollArea(blockLog)
	blockLogScroll.SetAutoscrollToBottom(true)

	blockBox := tui.NewVBox(blockLogScroll)
	blockBox.SetBorder(true)

	transactionLog = tui.NewVBox()

	transactionLogScroll := tui.NewScrollArea(transactionLog)
	transactionLogScroll.SetAutoscrollToBottom(true)

	transactionBox := tui.NewVBox(transactionLogScroll)
	transactionBox.SetBorder(true)

	blockBox.SetTitle("Block Stream")
	transactionBox.SetTitle("Transaction Stream")

	logsBox := tui.NewHBox(blockBox, transactionBox)

	input := tui.NewEntry()
	input.SetFocused(false)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	mainLogs := tui.NewVBox(logsBox, inputBox)
	mainLogs.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(sidebar, mainLogs)

	var err error
	ui, err = tui.New(root)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	client := p2p.Client{
		Handlers: []p2p.Handler{
			UILoggerHandler,
		},
	}

	err = client.Dial(":9876", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

var UILoggerHandler = func(processable p2p.PostProcessable) {

	p2pMsg := processable.P2PMessage
	switch p2pMsg.GetType() {

	case eos.SignedBlockMessageType:

		m := p2pMsg.(*eos.SignedBlockMessage)
		blockLog.Append(tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%s", m.Producer)),
			tui.NewPadder(1, 0, tui.NewLabel(m.Timestamp.Format("15:04:05"))),
			tui.NewSpacer(),
		))
		break

	case eos.SignedBlockSummaryMessageType:

		ui.Update(func() {
			m := p2pMsg.(*eos.SignedBlockSummaryMessage)
			blockLog.Append(tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("Block [%d]", m.BlockNumber())),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("By [%s]", m.Producer))),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("at [%s]", m.Timestamp.Format("15:04:05")))),
				tui.NewSpacer(),
			))
		})
		break
	case eos.PackedTransactionMessageType:

		ui.Update(func() {
			//m := p2pMsg.(*eos.PackedTransactionMessage)
			//transactionLog.Append(tui.NewHBox(
			//	tui.NewLabel("PackedTransactionMessage"),
			//	tui.NewSpacer(),
			//))
		})
		break
	case eos.SignedTransactionMessageType:

		ui.Update(func() {
			//m := p2pMsg.(*eos.PackedTransactionMessage)
			//transactionLog.Append(tui.NewHBox(
			//	tui.NewLabel("SignedTransactionMessageType"),
			//	tui.NewSpacer(),
			//))
		})
		break
	}

}
