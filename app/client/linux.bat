echo "set goos"
set GOOS=linux
echo "set GOPACH"
set GOPACH=amd64
echo "go build -o client client.go"
go build -o client client.go
