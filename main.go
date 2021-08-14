package main

import (
	"fmt"
	"github.com/mlsorensen/grafikeye/pkg/serial"
	"time"
)

func main() {
	mon := serial.QSESession{SerialPort: "/dev/ttyUSB0", BaudRate: 115200}
	err := mon.StartMonitor(handleMessage)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 5)
		fields := []string{"70", "3"}
		cmd := serial.QSCommand{serial.OperationExecute, serial.TypeDevice, serial.GrafikEye, fields}
		err = mon.Send(cmd)
		if err != nil {
			panic(err)
		}
	}
}

func handleMessage(command serial.QSCommand) {
	fmt.Printf("got this: %v\n", command)
}
