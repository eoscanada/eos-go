package p2p

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
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

	zlog.Info("Initiating proxy",
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
		zlog.Error("Started proxy error",
			zap.String("peer1", remoteAddress),
			zap.String("peer2", r.destinationPeerAddress),
			zap.Error(err))

		destinationPeer.connection.Close()
		remotePeer.connection.Close()

		zlog.Warn("Closing connection",
			zap.String("peer1", remoteAddress),
			zap.String("peer2", r.destinationPeerAddress))
		break
	case err := <-errorChannel:
		zlog.Error("Proxy error between %s and %s : %s",
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
			return errors.Wrapf(err, "peer init: listening %s", r.listeningAddress)
		}

		zlog.Info("Accepting connection", zap.String("listen", r.listeningAddress))

		for {
			conn, err := ln.Accept()
			if err != nil {
				zlog.Error("lost listening connection", zap.Error(err))
				break
			}
			zlog.Info("Connected to", zap.Stringer("remote", conn.RemoteAddr()))
			go r.startProxy(conn)
		}
	}
}
