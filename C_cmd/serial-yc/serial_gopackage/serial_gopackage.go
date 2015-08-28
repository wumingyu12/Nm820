package main

import (
	"github.com/huin/goserial"
	"log"
)

func main() {
	c := &goserial.Config{Name: "/dev/ttySAC1", Baud: 115200}
	s, err := goserial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test,hello"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2)
	for {
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buf[:n]))
	}

}
