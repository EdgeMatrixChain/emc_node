##setup_mac.sh
FILE=./emc
if test -d "$FILE"; then
    echo "$FILE exists."
else
    echo "downloading emc_mac.zip..."
    curl -o emc_linux_64.tgz https://edgematrix.pro/installer/emc_mac.zip
    echo "extracting zip..."
    unzip emc_mac.zip
    echo "init emc..."
    cd emc
    ./edge-matrix secrets init --data-dir edge_data
    echo "init complete"
fi