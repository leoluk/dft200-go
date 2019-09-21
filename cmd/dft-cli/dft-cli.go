package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/muka/go-bluetooth/api"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const (
	// rxChar receives notifications from the device
	rxCharUUID = "0000fff1-0000-1000-8000-00805f9b34fb"
	// txChar sends commands to the device
	txCharUUID = "0000fff2-0000-1000-8000-00805f9b34fb"
)

const (
	COMMAND_STATUS    uint8 = 0x00
	COMMAND_START     uint8 = 0x01
	COMMAND_STOP      uint8 = 0x02
	COMMAND_SET_SPEED uint8 = 0x03
	COMMAND_PAUSE     uint8 = 0x04
)

func checksum(b []byte) uint8 {
	var total uint8

	for _, n := range b {
		total += n
	}

	return total
}

func logToFile(w io.Writer, msg string) {
	timestamp := time.Now().Format(time.RFC1123)
	if _, err := fmt.Fprintf(w, "%s\t%s\n", timestamp, msg); err != nil {
		panic(err)
	}
}

func main() {
	addr := flag.String("addr", "", "Bluetooth device address")
	start := flag.Bool("start", false, "Start treadmill")
	stop := flag.Bool("stop", false, "Stop treadmill")
	pause := flag.Bool("pause", false, "Pause treadmill")
	speed := flag.Uint("speed", 0, "Set treadmill speed")
	logfile := flag.String("logfile", "", "Log actions to logfile")

	flag.Parse()

	if *addr == "" {
		log.Fatal("-addr is required")
	}

	dev, err := api.GetDeviceByAddress(*addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("device (dev): %v", dev)

	if !dev.IsConnected() {
		log.Infof("Connecting to device")
		err = dev.Connect()
	}

	if err != nil {
		log.Fatal(err)
	}

	var l *os.File
	if *logfile != "" {
		l, err = os.OpenFile(*logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
	}

	rxChar, err := dev.GetCharByUUID(rxCharUUID)
	if err != nil {
		log.Fatal(err)
	}

	txChar, err := dev.GetCharByUUID(txCharUUID)
	if err != nil {
		log.Fatal(err)
	}

	err = rxChar.StartNotify()
	if err != nil {
		log.Fatal(err)
	}

	payload := []byte("\xf0\xc3\x03")
	buffer := bytes.NewBuffer(payload)

	var (
		cmd   uint8
		value uint8
	)

	switch {
	case *start:
		cmd = COMMAND_START
		if l != nil {
			logToFile(l, "start")
		}
	case *stop:
		cmd = COMMAND_STOP
		if l != nil {
			logToFile(l, "stop")
		}
	case *speed != 0:
		cmd = COMMAND_SET_SPEED
		if l != nil {
			logToFile(l, fmt.Sprintf("speed=%d", *speed))
		}
	case *pause:
		cmd = COMMAND_PAUSE
		if l != nil {
			logToFile(l, "pause")
		}
	default:
		log.Fatal("No command specified")
	}

	err = binary.Write(buffer, binary.LittleEndian, cmd)
	if err != nil {
		panic(err)
	}

	switch {
	case *speed != 0:
		value = uint8(*speed)
	default:
		value = 0
	}

	err = binary.Write(buffer, binary.LittleEndian, value)
	if err != nil {
		panic(err)
	}

	err = binary.Write(buffer, binary.LittleEndian, uint8(0))
	if err != nil {
		panic(err)
	}

	err = binary.Write(buffer, binary.LittleEndian, checksum(buffer.Bytes()))
	if err != nil {
		panic(err)
	}

	err = txChar.WriteValue(buffer.Bytes(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Wrote payload %02x", buffer.Bytes())
}
