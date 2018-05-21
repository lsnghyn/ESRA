package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	uartPort = "/dev/ttyS0"
	baudRate = 9600
	timeOut  = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// readSerialGPIO will start reading bytes from the Serial
// Port and feed them thorugh the byte chan given on call.
func readSerialGPIO(data chan<- byte) {
	log.Println("Opening GPIO Serial")
	c := &serial.Config{Name: uartPort, Baud: baudRate, ReadTimeout: timeOut}
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
	NMEASentence := ""
	var b byte
	go readSerialGPIO(data)
	for {
		b = <-data
		NMEASentence += string(b)
		if b == byte(10) {
			fmt.Print(NMEASentence)
			NMEASentence = ""
		}
	}
}
