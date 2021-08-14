package main

import (
	"fmt"
	"github.com/mlsorensen/qse/pkg/serial"
)

func main() {
	mon := serial.QSESession{SerialPort: "/dev/ttyUSB0", BaudRate: 115200}
	err := mon.StartMonitor(handleMessage)
	if err != nil {
		panic(err)
	}

	select {}
}

func handleMessage(command serial.QSCommand) {
	fmt.Printf("got this: %v\n", command)
}
