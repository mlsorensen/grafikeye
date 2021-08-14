
default:
	GOOS=linux GOARCH=arm GOARM=7 go build main.go

test:
	go test ./...
