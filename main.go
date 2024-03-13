package main

import (
	"fmt"
	"net"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	proto = "udp"
	addr  = "192.168.1.1:7099"
)

var (
	cntr controller = controller{
		controlByte1:       128,
		controlByte2:       128,
		controlTurn:        128,
		controlAccelerator: 128,
	}
)

type controller struct {
	controlByte1       int
	controlByte2       int
	controlAccelerator int
	controlTurn        int
	isFastFly          bool
	isFastDrop         bool
	isEmergencyStop    bool
	isCircleTurnEnd    bool
	isNoHeadMode       bool
	isFastReturn       bool
	isUnLock           bool
	isGyroCorrection   bool
}

func (c *controller) access(i int) int {
	return (((c.controlByte1 ^ c.controlByte2) ^ c.controlAccelerator) ^ c.controlTurn) ^ (i & 255)
}

func main() {

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	conn := connection(proto, addr)

	ticker := time.NewTicker(300 * time.Millisecond)
	go processingPosition(ticker, conn)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyEsc:
			fmt.Println("Выход")
			ticker.Stop()
			return
		case keyboard.KeySpace:
			cntr.controlByte1 = 128
			cntr.controlByte2 = 128
			cntr.controlTurn = 128
			cntr.controlAccelerator = 128
		}

		switch char {
		case 'a':
			cntr.controlByte1 += 1
		case 'z':
			cntr.controlByte1 -= 1

		case 's':
			cntr.controlByte2 += 1
		case 'x':
			cntr.controlByte2 -= 1

		case 'd':
			cntr.controlAccelerator += 1
		case 'c':
			cntr.controlAccelerator -= 1

		case 'f':
			cntr.controlTurn += 1
		case 'v':
			cntr.controlTurn -= 1

		case 'q':
			cntr.isFastFly = !cntr.isFastFly
		case 'w':
			cntr.isFastDrop = !cntr.isFastDrop
		case 'e':
			cntr.isEmergencyStop = !cntr.isEmergencyStop
		case 'r':
			cntr.isCircleTurnEnd = !cntr.isCircleTurnEnd
		case 't':
			cntr.isNoHeadMode = !cntr.isNoHeadMode
		case 'y':
			cntr.isFastReturn = !cntr.isFastReturn
		case 'u':
			cntr.isUnLock = !cntr.isUnLock
		case 'i':
			cntr.isGyroCorrection = !cntr.isGyroCorrection
		}
	}
}

func processingPosition(ticker *time.Ticker, conn net.Conn) {
	for {
		select {
		case <-ticker.C:

			i := 0
			if cntr.isFastFly {
				i = 1
			}
			if cntr.isFastDrop {
				i += 2
			}
			if cntr.isEmergencyStop {
				i += 4
			}
			if cntr.isCircleTurnEnd {
				i += 8
			}
			if cntr.isNoHeadMode {
				i += 16
			}
			if cntr.isFastReturn || cntr.isUnLock {
				i += 32
			}
			if cntr.isGyroCorrection {
				i += 128
			}
			if cntr.controlTurn >= 104 && cntr.controlTurn <= 152 {
				// cntr.controlTurn = 128
			} else if cntr.controlTurn > 255 {
				cntr.controlTurn = 255
			} else if cntr.controlTurn < 1 {
				cntr.controlTurn = 1
			}
			if cntr.controlAccelerator == 1 {
				cntr.controlAccelerator = 0
			}
			if cntr.controlByte1 > 255 {
				cntr.controlByte1 = 255
			} else if cntr.controlByte1 < 1 {
				cntr.controlByte1 = 1
			}
			if cntr.controlByte2 > 255 {
				cntr.controlByte2 = 255
			} else if cntr.controlByte2 < 1 {
				cntr.controlByte2 = 1
			}

			fmt.Printf("tick %+v\n", cntr)

			if true {
				send(conn, []byte{1, 1}) // health check
				send(conn, []byte{
					3,
					102,
					byte(cntr.controlByte1),
					byte(cntr.controlByte2),
					byte(cntr.controlAccelerator),
					byte(cntr.controlTurn),
					byte(i),
					byte(cntr.access(i)),
					byte(153),
				})
			}
		}
	}
}

func connection(proto, addr string) net.Conn {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		panic(fmt.Sprint("err:", err))
	}
	return conn
}

func send(conn net.Conn, bts []byte) {
	fmt.Println("clnt send", bts)

	i, err := conn.Write(bts)
	if err != nil {
		panic(fmt.Sprint("err:", err))
	}
	fmt.Println("i ", i)

	buf := make([]byte, 20)
	d, err := conn.Read(buf)
	if err != nil {
		panic(fmt.Sprint("err:", err))
	}
	fmt.Println("clnt recv read", d)
	fmt.Println("clnt recv buf ", buf, "\n", string(buf))
}
