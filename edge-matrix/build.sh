LATEST_VERSION=0.18.1
LATEST_BUILD_VERSION=1
echo "LATEST_VERSION=$LATEST_VERSION"
echo "LATEST_BUILD_VERSION=$LATEST_BUILD_VERSION"
echo building linux-amd64...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/linux/emc/edge-matrix -ldflags="-s -w -X 'github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)' -X 'github.com/emc-protocol/edge-matrix/versioning.Build=$(echo $LATEST_BUILD_VERSION)'"
echo building windows-amd64...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../dist/windows/emc/edge-matrix.exe -ldflags="-s -w -X 'github.com/emc-protocol/edge-matrix/versioning.Version=$(echo $LATEST_VERSION)' -X 'github.com/emc-protocol/edge-matrix/versioning.Build=$(echo $LATEST_BUILD_VERSION)'"
echo completed.
