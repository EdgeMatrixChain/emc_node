echo "cp genesis.json to dist..."
cp ../genesis.json ../dist/windows/emc/
cp ../genesis.json ../dist/mac/emc/
cp ../genesis.json ../dist/mac_arm64/emc/
cp ../genesis.json ../dist/linux/emc/

#echo "package emc_windows_x64..."
#cd ../dist/windows/
#zip -r emc_windows_x64.zip emc
#mv emc_windows_x64.zip ..

#echo "package emc_mac..."
#cd ../dist/mac/
#zip -r emc_mac.zip emc
#mv emc_mac.zip ..
#
#echo "package emc_mac_arm64..."
#cd ../dist/mac_arm64/
#zip -r emc_mac_arm64.zip emc
#mv emc_mac_arm64.zip ..

echo "package emc_linux..."
cd ../dist/linux/
zip -r emc_linux_64.zip emc
mv emc_linux_64.zip ..

