package network

const (
// topicNameV1 = "edge/alive/0.2"
)

//func (s *Server) AddAddr(p peer.ID, addr multiaddr.Multiaddr) {
//	s.host.Peerstore().AddAddr(p, addr, peerstore.PermanentAddrTTL)
//
//	addrInfoString := addr.String() + "/p2p/" + p.String()
//	addrInfo, err := common.StringToAddrInfo(addrInfoString)
//	if err != nil {
//		s.logger.Error("AddAddr", "NodeId", p)
//		return
//	}
//	s.addPeerUpdateInfo(s.host.ID(), p, *addrInfo, false)
//}

//func (s *Server) addPeerUpdateInfo(from peer.ID, id peer.ID, addrInfo peer.AddrInfo, gossip bool) bool {
//	s.updatePeersLock.Lock()
//	defer s.updatePeersLock.Unlock()
//
//	s.logger.Info("addPeerUpdateInfo", "addrInfo", addrInfo.String(), "from", from.String(), "gossip", gossip)
//
//	updateInfo, updateInfoExists := s.updatePeers[id]
//
//	// Check if the peer update info is already initialized
//	needToPublish := false
//	if !updateInfoExists {
//		// Create a new record for the connection info
//		updateInfo = &PeerUpdateInfo{
//			Info:        addrInfo,
//			UpdateTime:  time.Now(),
//			From:        from.String(),
//			PublishTime: time.Now(),
//		}
//		s.updatePeers[id] = updateInfo
//		needToPublish = true
//	} else {
//		if s.updatePeers[id].PublishTime.After(s.updatePeers[id].UpdateTime.Add(peerstore.TempAddrTTL)) {
//			s.updatePeers[id].PublishTime = time.Now()
//			needToPublish = true
//		}
//		if s.updatePeers[id].PublishTime.After(s.updatePeers[id].UpdateTime.Add(peerstore.TempAddrTTL)) {
//			s.updatePeers[id].Info = addrInfo
//			s.updatePeers[id].From = from.String()
//			s.updatePeers[id].UpdateTime = time.Now()
//			needToPublish = true
//		}
//	}
//
//	// broadcast new addr
//	if !gossip && needToPublish && s.rtTopic != nil {
//		filteredPeers := make([]string, 0)
//		filteredPeers = append(filteredPeers, common.AddrInfoToString(&addrInfo))
//		s.rtTopic.Publish(&proto.PeerInfo{From: s.host.ID().String(), Nodes: filteredPeers})
//		s.logger.Debug("AddToPeerStore", "Publish", filteredPeers)
//	}
//
//	return false
//}

//func (s *Server) handlePeerAliveGossip(obj interface{}, from peer.ID) {
//	peerInfo, ok := obj.(*proto.PeerInfo)
//	if !ok {
//		s.logger.Error("failed to cast gossiped message to proto.PeerInfo")
//		return
//	}
//	if from.String() == s.host.ID().String() {
//		return
//	}
//
//	nodes := peerInfo.Nodes
//	for _, rawAddr := range nodes {
//		node, err := common.StringToAddrInfo(rawAddr)
//		if err != nil {
//			s.logger.Error("handlePeerAliveGossip", "err", fmt.Sprintf("failed to parse rawAddr %s: %w", rawAddr, err))
//			continue
//		}
//		s.host.Peerstore().AddAddr(node.ID, node.Addrs[0], peerstore.PermanentAddrTTL)
//		s.logger.Info("handlePeerAliveGossip", "from", from, "node", node.String())
//
//		s.addPeerUpdateInfo(from, node.ID, *node, true)
//	}
//}
//
//func (s *Server) StartPeerAliveGossip() error {
//	topic, err := s.NewTopic(topicNameV1, &proto.PeerInfo{})
//	if err != nil {
//		return err
//	}
//
//	if err := topic.Subscribe(s.handlePeerAliveGossip); err != nil {
//		return fmt.Errorf("unable to subscribe to gossip topic, %w", err)
//	}
//
//	s.rtTopic = topic
//	s.logger.Info("StartPeerAliveGossip")
//
//	return nil
//}
