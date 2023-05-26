##setup_mac_arm64.sh
FILE=./emc
if test -d "$FILE"; then
    echo "$FILE exists."
else
    echo "downloading emc_mac_arm64.zip..."
    curl -o emc_mac_arm64.zip https://edgematrix.pro/installer/emc_mac_arm64.zip
    echo "extracting zip..."
    unzip emc_mac_arm64.zip
    echo "init emc..."
    cd emc
    ./edge-matrix secrets init --data-dir edge_data
    echo "init complete"
fi