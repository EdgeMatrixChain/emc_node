package application

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Application struct {
	Name    string
	Tag     string
	Version string
	PeerID  peer.ID

	// ip4 addr string
	IpAddr string
	// ai app origin name
	AppOrigin string
	// ai model hash string
	ModelHash string
	// mac addr
	Mac string
	// memory info
	MemInfo string
	// cpu info
	CpuInfo string
	// gpu info
	GpuInfo string

	// app startup time
	StartupTime uint64
	// app uptime
	Uptime uint64
	// amount of slots currently occupying the app
	GuageHeight uint64
	// max limit
	GuageMax uint64
	// average e power value
	AveragePower float32
}

func (a *Application) Copy() *Application {
	newApp := &Application{
		Name:         a.Name,
		Tag:          a.Tag,
		Version:      a.Version,
		PeerID:       a.PeerID,
		StartupTime:  a.StartupTime,
		Uptime:       a.Uptime,
		GuageHeight:  a.GuageHeight,
		GuageMax:     a.GuageMax,
		IpAddr:       a.IpAddr,
		AppOrigin:    a.AppOrigin,
		Mac:          a.Mac,
		MemInfo:      a.MemInfo,
		CpuInfo:      a.CpuInfo,
		GpuInfo:      a.GpuInfo,
		ModelHash:    a.ModelHash,
		AveragePower: a.AveragePower,
	}

	return newApp
}
