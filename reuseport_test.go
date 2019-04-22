package libp2pquic

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-crypto"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

func TestPortReuse(t *testing.T) {
	prvKeyClient, _, _ := crypto.GenerateKeyPair(crypto.RSA, 2048)
	prvKeyServer, _, _ := crypto.GenerateKeyPair(crypto.RSA, 2048)

	tptClient, err := NewTransport(prvKeyClient, TransportOpt{})
	if err != nil {
		t.Fatal(err)
	}
	tptServer, err := NewTransport(prvKeyServer, TransportOpt{})
	if err != nil {
		t.Fatal(err)
	}

	runServer := func(tr tpt.Transport, multiaddr string) (ma.Multiaddr, <-chan tpt.Conn) {
		addrChan := make(chan ma.Multiaddr)
		connChan := make(chan tpt.Conn)
		go func() {
			addr, _ := ma.NewMultiaddr(multiaddr)
			ln, _ := tr.Listen(addr)
			addrChan <- ln.Multiaddr()
			conn, _ := ln.Accept()
			connChan <- conn
		}()
		return <-addrChan, connChan
	}

	serverMA, _ := runServer(tptServer, "/ip6/::1/udp/0/quic")
	_, err = tptClient.Dial(context.Background(), serverMA, "something")
	if err != nil {
		t.Fatal(err)
	}

}
