echo "cp genesis.json to dist..."
cp ../genesis.json ../dist/linux/emc/

echo "package emc_linux..."
cd ../dist/linux/
zip -r emc_linux_64.zip emc
mv emc_linux_64.zip ..