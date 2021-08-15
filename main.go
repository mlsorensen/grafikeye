package main

import (
	"fmt"
	. "github.com/mlsorensen/grafikeye/pkg/serial"
	"time"
)

func main() {
	session := QSESession{SerialPort: "/dev/ttyUSB0", BaudRate: 115200}
	err := session.StartMonitor(handleMessage)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 5)
		err := session.PressButton(ButtonScene1)
		if err != nil {
			panic(err)
		}
	}
}

func handleMessage(command QSCommand) {
	fmt.Printf("got this: %v\n", command)
}
