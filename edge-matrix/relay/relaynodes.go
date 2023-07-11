package relay

import (
	"errors"
	"sync/atomic"

	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	ErrNoRelaynodes  = errors.New("no relaynodes specified")
	ErrNoBootnodes   = errors.New("no bootnodes specified")
	ErrMinRelaynodes = errors.New("minimum 1 relaynode is required")
	ErrMinBootnodes  = errors.New("minimum 1 bootnode is required")
)

const (
	MinimumRelayNodes       int   = 1
	MinimumRelayConnections int64 = 1
)

type relaynodesWrapper struct {
	// relaynodeArr is the array that contains all the relaynode addresses
	relaynodeArr []*peer.AddrInfo

	// relaynodesMap is a map used for quick relaynode lookup
	relaynodesMap map[peer.ID]*peer.AddrInfo

	// relaynodeConnCount is an atomic value that keeps track
	// of the number of relaynode connections
	relaynodeConnCount int64
}

// isRelaynode checks if the node ID belongs to a set relaynode
func (bw *relaynodesWrapper) isRelaynode(nodeID peer.ID) bool {
	_, ok := bw.relaynodesMap[nodeID]

	return ok
}

// getRelaynodeConnCount loads the relaynode connection count [Thread safe]
func (bw *relaynodesWrapper) getRelaynodeConnCount() int64 {
	return atomic.LoadInt64(&bw.relaynodeConnCount)
}

// increaserelaynodeConnCount increases the relaynode connection count by delta [Thread safe]
func (bw *relaynodesWrapper) increaseRelaynodeConnCount(delta int64) {
	atomic.AddInt64(&bw.relaynodeConnCount, delta)
}

// getrelaynodes gets all the relaynodes
func (bw *relaynodesWrapper) getRelaynodes() []*peer.AddrInfo {
	return bw.relaynodeArr
}

// getrelaynodeCount returns the number of set relaynodes
func (bw *relaynodesWrapper) getRelaynodeCount() int {
	return len(bw.relaynodeArr)
}

// hasrelaynodes checks if any relaynodes are set [Thread safe]
func (bw *relaynodesWrapper) hasRelaynodes() bool {
	return bw.getRelaynodeCount() > 0
}
