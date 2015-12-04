echo "start complie main.go"
CC_FOR_TARGET=arm-linux-gcc GOOS=linux  GOARCH=arm GOARM=7 CGO_ENABLED=1 go build main.go