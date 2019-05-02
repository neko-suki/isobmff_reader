package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func parseMeta(file *os.File) uint32 {
	buf := make([]byte, 4)
	len, err := file.Read(buf)
	fmt.Println("len = ", len)
	if err != nil {
		log.Fatal(err)
	}
	version := uint8(buf[0])
	fmt.Println("version")
	fmt.Println(buf[1:3])
	flags := binary.BigEndian.Uint32(buf[1:3])
	fmt.Printf("version = %d, flags = %x\n", version, flags)

	_, err = file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	size := binary.BigEndian.Uint32(buf)

	_, err = file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	filetype := string(buf[0]) + string(buf[1]) + string(buf[2]) + string(buf[3])
	fmt.Println(filetype) // filetype

	size -= 8 // decrease unsigned int32 size and unsigned int32 filetype,
	_, err = file.Seek(int64(size), 1)
	if err != nil {
		log.Fatal(err)
	}
	return size
}

func main() {
	var fileName = flag.String("f", "default", "filename")
	flag.Parse()
	fmt.Println(*fileName)

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		buf := make([]byte, 4)
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		size := binary.BigEndian.Uint32(buf)

		_, err = file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		filetype := string(buf[0]) + string(buf[1]) + string(buf[2]) + string(buf[3])
		fmt.Println(filetype) // filetype

		size -= 8 // decrease unsigned int32 size and unsigned int32 filetype,
		switch filetype {
		case "meta":
			readSize := parseMeta(file)
			size -= readSize
		default:
		}
		_, err = file.Seek(int64(size), 1)
		if err != nil {
			log.Fatal(err)
		}
	}

}
