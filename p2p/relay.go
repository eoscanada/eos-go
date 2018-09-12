package p2p

import (
	"fmt"
	"net"
)

type Relay struct {
	listeningAddress       string
	destinationPeerAddress string
	handlers               []Handler
}

func NewRelay(listeningAddress string, destinationPeerAddress string) *Relay {
	return &Relay{
		listeningAddress:       listeningAddress,
		destinationPeerAddress: destinationPeerAddress,
	}
}

func (r *Relay) RegisterHandler(handler Handler) {

	r.handlers = append(r.handlers, handler)
}

func (r *Relay) startProxy(conn net.Conn) {

	remoteAddress := conn.RemoteAddr().String()

	fmt.Printf("Initiating proxy between %s and %s\n", remoteAddress, r.destinationPeerAddress)

	destinationPeer := NewOutgoingPeer(r.destinationPeerAddress, "eos-relay")

	errorChannel := make(chan error)

	destinationReadyChannel := destinationPeer.Connect(errorChannel)
	select {
	case <-destinationReadyChannel:
		remotePeer := newPeer(remoteAddress, fmt.Sprintf("agent-%s", remoteAddress), false)
		remotePeer.SetConnection(conn)
		proxy := NewProxy(remotePeer, destinationPeer)

		proxy.RegisterHandlers(r.handlers)

		err := proxy.Start("")
		fmt.Printf("Started proxy error between %s and %s : %s\n", remoteAddress, r.destinationPeerAddress, err)
		destinationPeer.connection.Close()
		remotePeer.connection.Close()
		fmt.Printf("Closing connection between %s and %s\n", remoteAddress, r.destinationPeerAddress)
		break
	case err := <-errorChannel:
		fmt.Printf("Proxy error between %s and %s : %s\n", conn.RemoteAddr(), r.destinationPeerAddress, err)
		break
	}
}

func (r *Relay) Start() error {

	for {
		ln, err := net.Listen("tcp", r.listeningAddress)
		if err != nil {
			return fmt.Errorf("peer init: listening %s: %s", r.listeningAddress, err)
		}

		fmt.Println("Accepting connection on:\n", r.listeningAddress)

		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("lost listening connection with: %s\n", err)
				break
			}
			fmt.Println("Connected to:", conn.RemoteAddr())
			go r.startProxy(conn)
		}
	}

	return nil
}
