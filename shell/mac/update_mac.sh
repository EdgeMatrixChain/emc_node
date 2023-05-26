FILE=./emc
if test -d "$FILE"; then
    echo "downloading emc_mac.zip..."
    mkdir -p tmp
    mkdir -p backup
    cd tmp
    curl -o emc_mac.zip https://edgematrix.pro/installer/emc_mac.zip
    echo "extracting zip..."
    unzip emc_mac.zip
    echo "backup emc..."
    mv ../emc/edge-matrix ../backup/
    mv ../emc/setup.sh ../backup/
    mv ../emc/start.sh ../backup/
    echo "update emc..."
    cp emc/edge-matrix ../emc
    cp emc/setup.sh ../emc
    cp emc/start.sh ../emc
    echo "clean up..."
    cd ..
    rm -rf tmp
    echo "update complete."
else
    echo "Can not do update for emc , $FILE is not exists."
fi
