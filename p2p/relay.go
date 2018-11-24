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

	logger.Infof("Initiating proxy between %s and %s",
		remoteAddress, r.destinationPeerAddress)

	destinationPeer := NewOutgoingPeer(r.destinationPeerAddress, "eos-relay", nil)

	errorChannel := make(chan error)

	destinationReadyChannel := destinationPeer.Connect(errorChannel)
	select {
	case <-destinationReadyChannel:
		remotePeer := newPeer(remoteAddress, fmt.Sprintf("agent-%s", remoteAddress), false, nil)
		remotePeer.SetConnection(conn)
		proxy := NewProxy(remotePeer, destinationPeer)

		proxy.RegisterHandlers(r.handlers)

		err := proxy.Start()
		logger.Errorf("Started proxy error between %s and %s : %s",
			remoteAddress, r.destinationPeerAddress, err)

		destinationPeer.connection.Close()
		remotePeer.connection.Close()

		logger.Warnf("Closing connection between %s and %s",
			remoteAddress, r.destinationPeerAddress)
		break
	case err := <-errorChannel:
		logger.Errorf("Proxy error between %s and %s : %s",
			conn.RemoteAddr(), r.destinationPeerAddress, err)
		break
	}
}

func (r *Relay) Start() error {

	for {
		ln, err := net.Listen("tcp", r.listeningAddress)
		if err != nil {
			return fmt.Errorf("peer init: listening %s: %s", r.listeningAddress, err)
		}

		logger.Infof("Accepting connection on:", r.listeningAddress)

		for {
			conn, err := ln.Accept()
			if err != nil {
				logger.Errorf("lost listening connection with: %s", err)
				break
			}
			logger.Info("Connected to:", conn.RemoteAddr())
			go r.startProxy(conn)
		}
	}

	return nil
}
