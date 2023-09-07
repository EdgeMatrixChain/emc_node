package application

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"math/big"
	"sync"
)

type AppPeer struct {
	// identifier
	ID string
	// name
	Name string
	// relay string
	Relay string
	// addr string
	Addr string
	// app origin name string
	AppOrigin string

	// ai model hash string
	ModelHash string
	// mac addr
	Mac string
	// memory info
	MemInfo string
	// cpu info
	CpuInfo string

	// peer's distance
	Distance *big.Int
	// app startup time
	Starup_time uint64
	// app uptime
	Uptime uint64
	// amount of slots currently occupying the app
	Guage_height uint64
	// max limit
	Guage_max uint64
	// average e power value
	AveragePower float32
	//gpu info
	GpuInfo string
	// version
	Version string
}

func (p *AppPeer) IsBetter(t *AppPeer) bool {
	if p.Guage_height != t.Guage_height {
		return p.Guage_height < t.Guage_height
	}

	return p.Distance.Cmp(t.Distance) < 0
}

type PeerMap struct {
	sync.Map
}

func NewPeerMap(peers []*AppPeer) *PeerMap {
	peerMap := new(PeerMap)

	peerMap.Put(peers...)

	return peerMap
}

func (m *PeerMap) Put(peers ...*AppPeer) {
	for _, peer := range peers {
		m.Store(peer.ID, peer)
	}
}

// Remove removes a peer from heap if it exists
func (m *PeerMap) Remove(peerID peer.ID) {
	m.Delete(peerID.String())
}

func (m *PeerMap) Get(id string) *AppPeer {
	value, ok := m.Load(id)
	if ok {
		return value.(*AppPeer)
	}
	return nil
}

// BestPeer returns the top of heap
func (m *PeerMap) BestPeer(skipMap map[string]bool) *AppPeer {
	var bestPeer *AppPeer

	m.Range(func(key, value interface{}) bool {
		peer, _ := value.(*AppPeer)

		if skipMap != nil && skipMap[peer.ID] {
			return true
		}

		if bestPeer == nil || peer.IsBetter(bestPeer) {
			bestPeer = peer
		}

		return true
	})

	return bestPeer
}
