FILE=./genesis.json
if test -f "$FILE"; then
  mv genesis.json genesis.json.backup
fi
curl -o genesis.json https://edgematrix.pro/installer/genesis.json
