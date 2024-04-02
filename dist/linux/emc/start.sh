./auto_select_rpc.sh
nohup ./edge-matrix server --chain genesis.json --data-dir edge_data  --grpc-address 0.0.0.0:50000 --libp2p 0.0.0.0:50001 --jsonrpc 0.0.0.0:50002 --base-libp2p 0.0.0.0:50003 --running-mode edge --app-name ComputingNode --log-to node.log --relay-on >edge_nohup.out 2>&1 &
