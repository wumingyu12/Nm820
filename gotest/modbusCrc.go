package main

import (
	"fmt"
)

func main() {
	cmd1 := []byte{0x00, 0x03, 0x00, 0x01, 0x00, 0x02}

	result_uint16 := Crc(cmd1) //uint16 格式,得到crc校验

	fmt.Printf("%x", result_uint16)
}

func Crc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}
