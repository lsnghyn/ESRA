package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarm/serial"
)

const (
	serialPort = "COM8" //"/dev/ttyS0"
	baudRate   = 57600
	timeOut    = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// sendSerialData will start reading a string from the
// input given by the user and will write to the Serial
// Port. ToDo: pass port s into function
func sendSerialData(data <-chan string) {
	log.Printf("Opening Serial Port: %v", serialPort)
	c := &serial.Config{Name: serialPort, Baud: baudRate, ReadTimeout: timeOut}
	s, err := serial.OpenPort(c)
	check(err)

	var str string
	for {
		str = <-data
		_, err := s.Write([]byte(str))
		check(err)
	}
}

func main() {
	data := make(chan string)
	reader := bufio.NewReader(os.Stdin)
	go sendSerialData(data)
	time.Sleep(time.Second * 2)
	for {
		// ToDo: get read data through bluetooth or wifi
		fmt.Print(">>>")
		sentence, err := reader.ReadString('\n')
		check(err)
		log.Printf("Sending Data: %v", sentence)
		data <- sentence
	}
}
