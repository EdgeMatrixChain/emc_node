# EMC Node

![](https://edgematrix.pro/_nuxt/logo.661b4f8c.png)

Beyond #ICP Layer2, serving as an entry of Computing power and Web3 in AI era

## Build
Execute the following command to compile the EMC node for linux_x64, windows_x64, mac(intel), mac_arm64(m1/m2)

```shell
cd ./emc_node/edge-matrix
sh build.sh
```

## Initial node
Execute the following command to init a EMC node.

```shell
cd ../dist/{linux|mac|mac_arm64}/emc
./edge-matrix secrets init --data-dir edge_data 
```
## Run
Execute the following command to run a EMC node.
Command with "--miner-canister nk6pr-3qaaa-aaaam-abnrq-cai" to works with the Testnet miner canister
```shell
./edge-matrix server --chain genesis.json --data-dir edge_data  --grpc-address 0.0.0.0:50000 --libp2p 0.0.0.0:50001 --jsonrpc 0.0.0.0:50002 --miner-canister nk6pr-3qaaa-aaaam-abnrq-cai 
```

## Basic Usage
Execute the following command to get help.
```shell
./edge-matrix help
```

## Java_sdk
https://github.com/EMCprotocol/emc_java_sdk

## Js-monorepo
https://github.com/EMCProtocol-dev/edgematrixjs-monorepo

## Sample
https://6tq33-2iaaa-aaaap-qbhpa-cai.icp0.io/

## Computing Node Test Tools
https://57hlm-riaaa-aaaap-qbhfa-cai.icp0.io

## Tutorials
For tutorials, check https://edgematrix.pro/start
