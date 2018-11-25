package p2p

import (
	"fmt"
	"net"

	"go.uber.org/zap"
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

	p2pLog.Info("Initiating proxy",
		zap.String("peer1", remoteAddress),
		zap.String("peer2", r.destinationPeerAddress))

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
		p2pLog.Error("Started proxy error",
			zap.String("peer1", remoteAddress),
			zap.String("peer2", r.destinationPeerAddress),
			zap.Error(err))

		destinationPeer.connection.Close()
		remotePeer.connection.Close()

		p2pLog.Warn("Closing connection",
			zap.String("peer1", remoteAddress),
			zap.String("peer2", r.destinationPeerAddress))
		break
	case err := <-errorChannel:
		p2pLog.Error("Proxy error between %s and %s : %s",
			zap.Stringer("peer1", conn.RemoteAddr()),
			zap.String("peer2", r.destinationPeerAddress),
			zap.Error(err))
		break
	}
}

func (r *Relay) Start() error {

	for {
		ln, err := net.Listen("tcp", r.listeningAddress)
		if err != nil {
			return fmt.Errorf("peer init: listening %s: %s", r.listeningAddress, err)
		}

		p2pLog.Info("Accepting connection", zap.String("listen", r.listeningAddress))

		for {
			conn, err := ln.Accept()
			if err != nil {
				logErr("lost listening connection", err)
				break
			}
			p2pLog.Info("Connected to", zap.Stringer("remote", conn.RemoteAddr()))
			go r.startProxy(conn)
		}
	}

	return nil
}
