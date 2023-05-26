FILE=./emc
if test -d "$FILE"; then
    echo "downloading emc_linux_64.tgz..."
    mkdir -p tmp
    mkdir -p backup
    cd tmp
    curl -o emc_linux_64.tgz https://edgematrix.pro/installer/emc_linux_64.tgz
    echo "extracting tgz..."
    tar xzf emc_linux_64.tgz
    echo "backup emc..."
    mv ../emc/edge-matrix ../backup/
    mv ../emc/start.sh ../backup/
    mv ../emc/genesis.json ../backup/
    echo "update emc..."
    cp emc/edge-matrix ../emc
    cp emc/start.sh ../emc
    cp emc/genesis.json ../emc
    echo "clean up..."
    cd ..
    rm -rf tmp
    echo "update complete."
else
    echo "Can not do update for emc , $FILE is not exists."
fi