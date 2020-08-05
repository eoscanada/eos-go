package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/blockslog"
	"github.com/eoscanada/eos-go/p2p"
	"github.com/pkg/errors"
)

var peer = flag.String("peer", "localhost:9876", "peer to connect to")
var blocksLog = flag.String("blocks-log-path", "blocks/blocks.log", "Path to a valid blocks.log file")
var showLog = flag.Bool("v", false, "show detail log")

func main() {
	flag.Parse()

	if *showLog {
		p2p.EnableP2PLogging()
	}
	defer p2p.SyncLogger()

	blkReader := blockslog.NewReader(*blocksLog)
	defer blkReader.Close()
	if err := blkReader.ReadHeader(); err != nil {
		log.Fatal("read header", err)
	}

	// firstBlk, err := blkReader.Next()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err := blkReader.Last(); err != nil {
		log.Fatal("last", err)
	}

	lastBlk, _, err := blkReader.Prev()
	if err != nil {
		log.Fatal("prev", err)
	}

	cID, err := hex.DecodeString(blkReader.ChainID)
	if err != nil {
		log.Fatal("decode chain id", err)
	}

	lastBlockID, err := lastBlk.BlockID()
	if err != nil {
		log.Fatal("block id compute", err)
	}

	handshake := &p2p.HandshakeInfo{
		ChainID:                  cID,
		HeadBlockNum:             lastBlk.BlockNumber(),
		HeadBlockID:              lastBlockID,
		HeadBlockTime:            lastBlk.Timestamp.Time,
		LastIrreversibleBlockNum: lastBlk.BlockNumber(),
		LastIrreversibleBlockID:  lastBlockID,
	}

	fmt.Println("Connect to ", *peer, " with Chain ID:", blkReader.ChainID)
	client := NewClient(
		p2p.NewOutgoingPeer(*peer, "blockfeeder", handshake),
		handshake,
		blkReader,
	)

	client.Start()
}

type Client struct {
	peer        *p2p.Peer
	handshake   *p2p.HandshakeInfo
	readTimeout time.Duration
	blkReader   *blockslog.Reader
}

func NewClient(peer *p2p.Peer, handshake *p2p.HandshakeInfo, blkReader *blockslog.Reader) *Client {
	client := &Client{
		peer:      peer,
		handshake: handshake,
		blkReader: blkReader,
	}
	return client
}

func (c *Client) SetReadTimeout(readTimeout time.Duration) {
	c.readTimeout = readTimeout
}

func (c *Client) read(peer *p2p.Peer, errChannel chan error) {
	for {
		packet, err := peer.Read()
		if err != nil {
			errChannel <- fmt.Errorf("read message from %s: %w", peer.Address, err)
			break
		}

		//envelope := p2p.NewEnvelope(peer, peer, packet)

		fmt.Printf("Incoming message: %T %v\n", packet.P2PMessage, packet.P2PMessage)
		switch m := packet.P2PMessage.(type) {
		case *eos.GoAwayMessage:
			errChannel <- errors.Wrapf(err, "GoAwayMessage reason %s", m.Reason)
		case *eos.HandshakeMessage:
			///fmt.Println("MAMA", m.LastIrreversibleBlockNum, c.handshake.LastIrreversibleBlockNum)
			if m.LastIrreversibleBlockNum < c.handshake.LastIrreversibleBlockNum {
				fmt.Println("Writing notice message")
				err := c.peer.WriteP2PMessage(&eos.NoticeMessage{
					KnownTrx: eos.OrderedSelectIDs{
						Mode: 0, /* mode == none */
					},
					KnownBlocks: eos.OrderedSelectIDs{
						Mode:    3, /* mode == normal */
						Pending: c.handshake.HeadBlockNum,
						IDs:     []eos.Checksum256{c.handshake.HeadBlockID},
					},
				})
				if err != nil {
					errChannel <- err
				}

				// Start PUSHING blocks! from their `LastIrreversibleBlockNum`
				if err := c.pushBlocks(m.LastIrreversibleBlockNum); err != nil {
					errChannel <- err
				}

			}
		case *eos.NoticeMessage:
		case *eos.SignedBlock:
		default:
		}
	}
}

func (c *Client) pushBlocks(fromBlockNum uint32) error {
	c.blkReader.First()
	for {
		blk, rawBytes, err := c.blkReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		blkNum := blk.BlockNumber()
		if blkNum < fromBlockNum {
			fmt.Println("Skipping block", blkNum)
			continue
		}

		for i := 0; i < 100; i++ {
			fmt.Println("Writing block", blkNum)
			err = c.peer.WritePacket(&eos.Packet{
				Type:    eos.SignedBlockType,
				Payload: rawBytes,
			})
			if err != nil {
				return err
			}
		}
	}
}

func (c *Client) Start() error {
	errorChannel := make(chan error, 1)
	readyChannel := c.peer.Connect(errorChannel)

	for {
		select {
		case <-readyChannel:
			go c.read(c.peer, errorChannel)

			err := c.peer.SendHandshake(c.handshake)
			if err != nil {
				return fmt.Errorf("start: send handshake: %w", err)
			}

		case err := <-errorChannel:
			return fmt.Errorf("start failed: %w", err)
		}
	}
}
