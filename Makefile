
default:
	GOOS=linux GOARCH=arm GOARM=7 go build -o grafikeye main.go

test:
	go test ./...
