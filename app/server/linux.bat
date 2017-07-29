echo "set goos"
set GOOS=linux
echo "set GOPACH"
set GOPACH=amd64
echo "go build -o -x umbrella server.go"
go build -o umbrella server.go
