package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/adrianmo/go-nmea"
	"github.com/tarm/serial"
)

const (
	serialPort = "COM4"
	baudRate   = 57600
	timeOut    = time.Second * 60
)

// check is a simple function that checks an error and
// fatally logs on error.
func check(err error) {
	if err == io.EOF {
		return
	} else if err != nil {
		log.Panic(err)
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
	data := make(chan byte)
	sentence := ""
	var b byte
	go readSerial(data)
	for {
		b = <-data
		sentence += string(b)
		if b == byte(10) {
			_, err = file.WriteString(sentence)
			s, err := nmea.Parse(sentence)
			if err != nil {
				log.Fatal(err)
			}
			m := s.(nmea.GPRMC)
			log.Printf("Raw sentence: %v\n", m)
			log.Printf("Time: %s\n", m.Time)
			log.Printf("Validity: %s\n", m.Validity)
			log.Printf("Latitude GPS: %s\n", nmea.FormatGPS(m.Latitude))
			log.Printf("Latitude DMS: %s\n", nmea.FormatDMS(m.Latitude))
			log.Printf("Longitude GPS: %s\n", nmea.FormatGPS(m.Longitude))
			log.Printf("Longitude DMS: %s\n", nmea.FormatDMS(m.Longitude))
			log.Printf("Speed: %f\n", m.Speed)
			log.Printf("Course: %f\n", m.Course)
			log.Printf("Date: %s\n", m.Date)
			log.Printf("Variation: %f\n", m.Variation)
			sentence = ""
		}
	}
}
