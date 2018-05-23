package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	serialPort = "/dev/ttyUSB0"
	baudRate   = 57600
	timeOut    = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// readSerial will start reading bytes from the Serial
// Port and feed them thorugh the byte chan given on call.
func readSerial(data chan<- byte) {
	log.Printf("Opening Serial Port: %v\n", serialPort)
	c := &serial.Config{Name: serialPort, Baud: baudRate, ReadTimeout: timeOut}
	s, err := serial.OpenPort(c)
	check(err)
	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		check(err)
		for _, b := range buf[:n] {
			data <- b
		}
	}
}

func main() {
	data := make(chan byte)
	sentence := ""
	var b byte
	go readSerial(data)
	for {
		b = <-data
		sentence += string(b)
		if b == byte(10) {
			fmt.Print(sentence)
			sentence = ""
		}
	}
}
