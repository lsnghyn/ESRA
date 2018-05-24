package main

import (
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	uartPort   = "/dev/ttyS0"
	serialPort = "/dev/ttyUSB0"
	baudRate   = 9600
	timeOut    = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err != nil {
		log.Panic(err)
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

// sendSerialData will start reading a string from the
// input given by the user and will write to the Serial
// Port. ToDo: pass port s into function
func sendSerialData(data <-chan string) {
	c := &serial.Config{Name: serialPort, Baud: 57600, ReadTimeout: timeOut}
	s, err := serial.OpenPort(c)
	check(err)

	for {
		_, err := s.Write([]byte(<-data))
		check(err)
	}
}

func main() {
	dataIn := make(chan byte)
	dataOut := make(chan string)
	sentence := ""
	var b byte
	go readSerialGPIO(dataIn)
	go sendSerialData(dataOut)
	for {
		b = <-dataIn
		sentence += string(b)
		if b == byte(10) {
			log.Print(sentence)
			dataOut <- sentence
			sentence = ""
		}
	}
}
