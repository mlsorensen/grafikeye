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
	BaudRate   int
	openPort   *serial.Port
}

type QSCommand struct {
	Operation     byte
	Type          string
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
// TODO: handle closing session
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
				time.Sleep(time.Second * 10)
				continue
			}

			for i := 0; i < num; i++ {
				if buf[i] == '\r' {
					// log string and error
					msg := strings.TrimSpace(string(line))
					line = []byte{}
					cmd, err := parseQSEMessage(msg)
					if err != nil {
						fmt.Printf("Error parsing line '%s': %v\n", msg, err)
						continue
					}
					callback(cmd)
				} else {
					line = append(line, buf[i])
				}
			}
		}
	}()

	return nil
}

func (q *QSESession) Send(cmd QSCommand) error {
	if q.openPort == nil {
		err := q.NewSession()
		if err != nil {
			return err
		}
	}

	_, err := q.openPort.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (q *QSESession) PressButton(button string) error {
	pressFields := []string{button, ActionPress}
	releaseFields := []string{button, ActionRelease}
	cmd := QSCommand{OperationExecute, TypeDevice, GrafikEye, pressFields}
	err := q.Send(cmd)
	if err != nil {
		return err
	}
	cmd.CommandFields = releaseFields
	return q.Send(cmd)
}

func (c *QSCommand) Bytes() []byte {
	opType := string(c.Operation) + c.Type
	cmdFields := strings.Join(c.CommandFields, ",")
	str := fmt.Sprintf("%s,%s,%s\r\n", opType, c.IntegrationId, cmdFields)
	return []byte(str)
}

func parseQSEMessage(line string) (cmd QSCommand, err error) {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, "QSE>")
	parts := strings.Split(line, ",")

	if len(parts) == 0 {
		err = fmt.Errorf("unable to parse line '%s'", line)
		return
	}

	operation := parts[0][0]
	operationMatcher := fmt.Sprintf("^[%c%c%c]$", OperationMonitor, OperationExecute, OperationQuery)
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
