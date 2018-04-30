package main

import (
	"fmt"
	"log"

	"bytes"
	"net/url"

	"flag"

	"encoding/hex"

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

var apiAddr = flag.String("api-addr", "http://localhost:8888", "RPC endpoint of the nodeos instance")
var p2pAddr = flag.String("p2p-addr", "localhost:8902", "P2P socket connection")
var signingKey = flag.String("signing-key", "", "Key to sign transactions we're about to blast")
var chainID = flag.String("chain-id", "00000000000000000000000000000000", "Chain id")
var networkVersion = flag.Int("network-version", 25431, "Chain id")

func main() {

	flag.Parse()

	apiAddrURL, err := url.Parse(*apiAddr)
	if err != nil {
		log.Fatalln("could not parse --api-addr:", err)
	}

	//r, w, _ := os.Pipe()
	//os.Stdout = w
	//
	//go func() {
	//	scanner := bufio.NewScanner(r)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//
	//		ui.Update(func() {
	//			blockLog.Append(tui.NewHBox(
	//
	//				tui.NewLabel(line),
	//				tui.NewSpacer(),
	//			))
	//		})
	//	}
	//}()

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

	ui, err = tui.New(root)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	api := eos.New(apiAddrURL, bytes.Repeat([]byte{0}, 32))
	client := p2p.NewClient(*p2pAddr, api, *chainID, int16(*networkVersion))
	client.RegisterHandler(p2p.HandlerFunc(UILoggerHandler))
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

var UILoggerHandler = func(processable p2p.PostProcessable) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

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
				tui.NewLabel(fmt.Sprintf("[%d]", m.BlockNumber())),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("By [%s]", m.Producer))),
				tui.NewLabel(fmt.Sprintf("at [%s]", m.Timestamp.Format("15:04:05"))),
				tui.NewSpacer(),
			))
			for _, region := range m.Regions {
				for _, cycleSummary := range region.CyclesSummary {
					for _, cycle := range cycleSummary {
						for _, tx := range cycle.Transactions {
							blockLog.Append(tui.NewHBox(
								tui.NewLabel(fmt.Sprintf("  Tx [%s] status [%s]", tx.Status, hex.EncodeToString(tx.ID[:4]))),
							))
						}
					}
				}
			}
		})
		break
	case eos.PackedTransactionMessageType:

		m := p2pMsg.(*eos.PackedTransactionMessage)
		signedTx, err := m.UnPack()
		if err != nil {
			fmt.Println("PackedTransactionMessage: ", err)
			return
		}
		fmt.Println(signedTx)
		ui.Update(func() {

			transactionLog.Append(tui.NewHBox(
				//tui.NewLabel(fmt.Sprintf("[%s]", signedTx.Transaction.ID())),
				tui.NewLabel("PackedTransactionMessage !!!!!!!!!!!!!!!!!!!!!!!!!!"),
			))
		})
		break
	case eos.SignedTransactionMessageType:

		ui.Update(func() {
			//m := p2pMsg.(*eos.PackedTransactionMessage)
			transactionLog.Append(tui.NewHBox(
				tui.NewLabel("SignedTransactionMessageType"),
				tui.NewSpacer(),
			))
		})
		break
	}

}
