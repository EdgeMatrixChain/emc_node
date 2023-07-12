package application

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Application struct {
	Name    string
	Tag     string
	Version string
	PeerID  peer.ID

	//RelayInfo *relay.RelayPeerInfo

	// app startup time
	StartupTime uint64
	// app uptime
	Uptime uint64
	// amount of slots currently occupying the app
	GuageHeight uint64
	// max limit
	GuageMax uint64
}

func (a *Application) Copy() *Application {
	newApp := &Application{
		Name:        a.Name,
		Tag:         a.Tag,
		Version:     a.Version,
		PeerID:      a.PeerID,
		StartupTime: a.StartupTime,
		Uptime:      a.Uptime,
		GuageHeight: a.GuageHeight,
		GuageMax:    a.GuageMax,
		//RelayInfo:   a.RelayInfo,
	}

	return newApp
}
