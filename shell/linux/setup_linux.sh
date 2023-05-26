##setup_linux.sh
FILE=./emc
if test -d "$FILE"; then
    echo "$FILE exists."
else
    echo "downloading emc_linux_64.tgz..."
    curl -o emc_linux_64.tgz https://edgematrix.pro/installer/emc_linux_64.tgz
    echo "extracting tgz..."
    tar xzf emc_linux_64.tgz
    echo "init emc..."
    cd emc
    ./edge-matrix secrets init --data-dir edge_data
    echo "init complete"
fi