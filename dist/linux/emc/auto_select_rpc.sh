#!/bin/bash

content=$(curl -s https://install.edgematrix.pro/list)

min_latency=1000
fastest_ip=""
fastest_url=""
fastest_area=""

while IFS= read -r line; do
    ip_list=$(echo $line | awk '{print $3}' FS=" ")
    ip_name=$(echo $line | awk '{print $2}' FS=" ")
    ip_area=$(echo $line | awk '{print $1}' FS=" ")

    IFS=',' read -r -a ip_array <<< "$ip_list"

    for ip in "${ip_array[@]}"
    do
        latency=$(ping -c 5 $ip | grep 'avg' | awk -F'/' '{print $5}')
        echo "Ping $ip avg relay: $latency ms"

        if (( $(echo "$latency < $min_latency" | bc -l) ))
        then
            min_latency=$latency
            fastest_ip=$ip
            fastest_url=$ip_name
            fastest_area=$ip_area
        fi
    done
done <<< "$content"

echo "the fastest area is: $fastest_area, IP: $fastest_ip, relay time: $min_latency ms"

curl -o genesis.json -O https://install.edgematrix.pro/genesis.json.$fastest_area
nodeId=$(./edge-matrix secrets output --node-id --data-dir ./edge_data/)
curl --location --request POST 'https://openapi.emchub.ai/emchub/api/client/open/reportNodeRpc' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'nodeId='$nodeId \
    --data-urlencode 'rpcAddress='$fastest_url