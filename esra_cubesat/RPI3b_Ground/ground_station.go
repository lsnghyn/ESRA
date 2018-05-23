package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

const (
	serialPort = "/dev/ttyS0"
	baudRate   = 9600 //57600
	timeOut    = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err == io.EOF {
		return
	} else if err != nil {
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
	file, err := os.Create("gps_data.txt")
	check(err)
	defer func() {
		check(file.Close())
	}()
	w := bufio.NewWriter(file)
	data := make(chan byte)
	sentence := ""
	var b byte
	go readSerial(data)
	for {
		b = <-data
		sentence += string(b)
		if b == byte(10) {
			_, err = w.Write([]byte(sentence))
			log.Print(sentence)
			sentence = ""
		}
	}
}
