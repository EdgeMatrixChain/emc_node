LATEST_VERSION=0.9.16
echo building darwin-amd64...
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../dist/mac/emc/edge-matrix -ldflags="-s -w -X github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)"
echo building darwin-arm64...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../dist/mac_arm64/emc/edge-matrix -ldflags="-s -w -X github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)"
echo building linux-amd64...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/linux/emc/edge-matrix -ldflags="-s -w -X github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)"
echo building windows-amd64...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../dist/windows/emc/edge-matrix.exe -ldflags="-s -w -X github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)"
echo completed.