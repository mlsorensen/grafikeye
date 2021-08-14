package serial

import (
	"fmt"
	"github.com/tarm/serial"
	"regexp"
	"strings"
	"time"
)

type QSESession struct {
	SerialPort string
	BaudRate int
	openPort *serial.Port
}

type QSCommand struct {
	Operation byte
	Type string
	IntegrationId string
	CommandFields []string // using string because field content types are variable
}

func (q *QSESession) NewSession() error {
	c := &serial.Config{Name: q.SerialPort, Baud: q.BaudRate}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	q.openPort = s

	return nil
}


// StartMonitor begins reading the serial connection. When a complete
// QSE message is found, it calls the provided function.
// TODO: add cancel channel to break out of monitor loop
func (q *QSESession) StartMonitor(callback func(command QSCommand)) error {
	if q.openPort == nil {
		err := q.NewSession()
		if err != nil {
			return err
		}
	}

	var line []byte
	buf := make([]byte, 128)

	go func() {
		for {
			num, err := q.openPort.Read(buf)
			if err != nil {
				// log error, throttle retries to 10s
				time.Sleep(10)
				continue
			}

			for i := 0; i < num; i++ {
				if buf[i] == '\r' {
					// log string and error
					cmd, _ := parseQSEMessage(string(line))
					callback(cmd)
					line = []byte{}
				} else {
					line = append(line, buf[i])
				}
			}
		}
	}()

	return nil
}

func parseQSEMessage(line string) (cmd QSCommand, err error) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "QSE>")
	parts := strings.Split(line, ",")

	if len(parts) == 0 {
		err = fmt.Errorf("unable to parse line '%s'", line)
		return
	}

	operation := parts[0][0]
	operationMatcher := fmt.Sprintf("^[%c%c%c]$",OperationMonitor, OperationExecute, OperationQuery)
	validOperation, err := regexp.Match(operationMatcher, []byte{operation})
	if err != nil {
		err = fmt.Errorf("unable to parse line '%s'", line)
		return
	}
	if !validOperation {
		err = fmt.Errorf("invalid operation found '%c' in line '%s'", operation, line)
	}

	cmd.Operation = operation
	cmd.Type = parts[0][1:]
	cmd.IntegrationId = parts[1]
	cmd.CommandFields = parts[2:]

	return
}