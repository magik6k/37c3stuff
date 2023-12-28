package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// {"delay":0,"reporting_point":47490,"priority":0,"direction_request":0,"line":429,"run_number":42,"destination_number":943,"train_length":0}

	var in struct {
		Delay             int `json:"delay"`
		ReportingPoint    int `json:"reporting_point"`
		Priority          int `json:"priority"`
		DirectionRequest  int `json:"direction_request"`
		Line              int `json:"line"`
		RunNumber         int `json:"run_number"`
		DestinationNumber int `json:"destination_number"`
		TrainLength       int `json:"train_length"`
	}

	//data := os.Args[1]
	// data is readline from stdin
	/*data := ""
	b := bufio.NewReader(os.Stdin)
	db, err := b.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	data = string(db)*/

	var data string

	// connect to 45.158.40.165 1337, read 5 lines
	conn, err := net.Dial("tcp", "45.158.40.165:1337")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	b := bufio.NewReader(conn)
	var lines []string
	for i := 0; i < 3; i++ {
		line, err := b.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("line %d: %s\n", i, line)
		lines = append(lines, line)
	}

	fmt.Println(lines[0])
	data = lines[0]

	toTrim := "Here is your data for the R09.16 telegram: "
	data = strings.TrimPrefix(data, toTrim)

	err = json.Unmarshal([]byte(data), &in)
	if err != nil {
		panic(err)
	}

	/*
		111111 000000000 1 xxxx xxxx 1 xxxx xxxx 1 xxxx xxxx 1 xxxx xxxx 1 00000000 1 00000000

		11111100000000 1 10010001100000010100000000111000000100000000100000000

		00010001 00000110 MP MP 00000000 rounte_num

		MMMMMMMM MMMMMMMM PPHHLLLL LLLLLLLL KKKKKKKK ZZZZZZZZ ZZZZRLLL

		01000110 11110100 00000011 00001110 01010000 00011101 10100000

		01000110111101000000001100001110010100000001110110100000

		////

		b1    00010001 / ?0000
		b2    00000110 /
		b3 MP
		a1 MP
		a2 PR 00000000 // 11?? dir /
		a3 LN  line number
		a4 KN run_number 1010000
		a5 ZS dest 0000 0000
		a6 0000 0000


		we wann to output a binary packet

	*/

	/*
		// b1
			fmt.Println("00010001")
			fmt.Println("1") // byte sep
			// b2
			fmt.Println("00000110")
			fmt.Println("1") // byte sep
			// b3 / a1 MP
			fmt.Printf("%016b\n", in.ReportingPoint)
			fmt.Println("1") // byte sep
			// a2 PR 00000000 // 11?? dir /
			fmt.Printf("%02b\n", in.Priority)
			fmt.Printf("%02b\n", in.DirectionRequest)

			fmt.Printf("%04b\n", in.Line&0xF)
			fmt.Println("1") // byte sep
			// a3 LN  line number
			fmt.Printf("%08b\n", in.Line>>4)
			fmt.Println("1") // byte sep
			// a4 KN run_number
			fmt.Printf("%07b\n", in.RunNumber)
			fmt.Println("1") // byte sep
			// a5 ZN dest
			fmt.Printf("%08b\n", in.DestinationNumber&0xFF)
			fmt.Println("1") // byte sep
			// a6 ZN / ZL
			fmt.Printf("%04b\n", in.DestinationNumber>>8)
			fmt.Println("0")
			fmt.Printf("%03b\n", in.TrainLength)
			fmt.Println("1") // byte sep
	*/

	arg := func(char string) bool {
		a := ""
		if len(os.Args) > 1 {
			a = os.Args[1]
		}
		return strings.Contains(a, char)
	}

	bsync0 := arg("b")
	bsync1 := arg("f")
	startbit := arg("s")
	bytesep := arg("p")

	endianness := true

	var output string

	bsep := func() {
		if bytesep {
			fmt.Println("BSEP 1") // byte sep
			output += "1"
		}
	}

	var prefix string

	if bsync0 {
		fmt.Println("BSYNC 111111")
		output += "111111"
		prefix = "111111"
	}
	if bsync1 {
		fmt.Println("FSYNC 0000000000")
		output += "0000000000"
		prefix += "0000000000"
	}
	if startbit {
		fmt.Println("SBIT 1")
		output += "1"
		prefix += "1"
	}
	//                      M           T           MP          MP          PR            ln         kn          zn           zn2
	// m 11111100000000001 1001 0001 1 0000 0110 1 1010 0111 1 1000 1011 1 0000 0011 1 1100 1110 1 0001 1101 1 0000 0110 1 1001 0000 100000000100000000
	// t 11111100000000001 1000 1001 1 0110 0000 1 0000 1101 1 1101 1011 1 0011 0010 1 1111 0100 1 0011 1001 1 0101 0000 1 1101 0010 111010000100101011
	// e 11111100000000001 1000 1001 1 0110 0000 1 1101 0110 1 1011 1000 1 1000 1010 1 1011 1110 1 1110 0100 1 0100 0100 1 0100 0011 111101011110111010
	// b 11111100000000001 1000 1001 1 0110 0000 1 0000 1000 1 0100 0111 1 1000 0000 1 0111 1100 1 0111 0100 1 1101 0100 1 0000 0100 100000000100000000
	// a 111111000000000110001001101100000101101001111101111110000000110011101101001010111110000100001010100000000100000000

	if arg("1") {
		fmt.Println("B1b1 1, the :")
		output += "1"
	} else {
		fmt.Println("B1b1 0, the :")
		output += "0"
	}
	fmt.Println("B1 0010001 MMM TTTT")
	output += "0010001"
	bsep()

	// b2
	if arg("V") {
		fmt.Println("B2b1 1 ZV delay")
		output += "1"
	} else {
		fmt.Println("B2b1 0 ZV delay")
		output += "0"
	}
	fmt.Println("B2 0000110 WWW delay LLLL len")
	output += "0000110"
	bsep()
	// b3 / a1 MP

	//fmt.Printf("%016b\n", in.ReportingPoint)
	if endianness {
		fmt.Printf("B3 MP %08b MMMMMMMM\n", in.ReportingPoint>>8)
		output += fmt.Sprintf("%08b", in.ReportingPoint>>8)
		bsep()
		fmt.Printf("A1 MP %08b MMMMMMMM\n", in.ReportingPoint&0xFF)
		output += fmt.Sprintf("%08b", in.ReportingPoint&0xFF)
	} else {
		fmt.Printf("%08b\n", in.ReportingPoint&0xFF)
		output += fmt.Sprintf("%08b", in.ReportingPoint&0xFF)
		bsep()
		fmt.Printf("%08b\n", in.ReportingPoint>>8)
		output += fmt.Sprintf("%08b", in.ReportingPoint>>8)
	}

	bsep()
	// a2 PR 00000000 // 11?? dir /
	fmt.Printf("PR %02b prio\n", in.Priority)
	output += fmt.Sprintf("%02b", in.Priority)
	fmt.Printf("DR %02b dir\n", in.DirectionRequest)
	output += fmt.Sprintf("%02b", in.DirectionRequest)

	if endianness {
		fmt.Printf("LN %04b LLLL line\n", in.Line>>8)
		output += fmt.Sprintf("%04b", in.Line>>8)
		bsep()
		// a3 LN  line number
		fmt.Printf("LN %08b LLLLLLLL line\n", in.Line&0xFF)
		output += fmt.Sprintf("%08b", in.Line&0xFF)
	} else {
		fmt.Printf("%08b\n", in.Line>>4)
		output += fmt.Sprintf("%08b", in.Line>>4)
		bsep()
		// a3 LN  line number
		fmt.Printf("%04b\n", in.Line&0xF)
		output += fmt.Sprintf("%04b", in.Line&0xF)
	}

	bsep()
	// a4 KN run_number
	fmt.Printf("KN %08b KKKKKKKK run number\n", in.RunNumber)
	output += fmt.Sprintf("%08b", in.RunNumber)
	bsep()

	// a5 ZN dest
	if endianness {
		fmt.Printf("ZN %08b ZZZZZZZZ dest\n", in.DestinationNumber>>4)
		output += fmt.Sprintf("%08b", in.DestinationNumber>>4)
		bsep()
		// a6 ZN / ZL
		fmt.Printf("ZN %04b ZZZZ dest\n", in.DestinationNumber&0xF)
		output += fmt.Sprintf("%04b", in.DestinationNumber&0xF)
	} else {
		fmt.Printf("%04b\n", in.DestinationNumber&0xF)
		output += fmt.Sprintf("%04b", in.DestinationNumber&0xF)
		bsep()
		// a6 ZN / ZL
		fmt.Printf("%08b\n", in.DestinationNumber>>4) // todo makes no sense here, should be 4 bits
		output += fmt.Sprintf("%08b", in.DestinationNumber>>4)
	}
	fmt.Println("R 0")
	output += "0"

	fmt.Printf("ZL %03b\n", in.TrainLength)
	output += fmt.Sprintf("%03b", in.TrainLength)
	bsep()

	// 2 crc bytes
	fmt.Println("CRC 00000000")
	output += "00000000"
	bsep()
	fmt.Println("CRC 00000000")
	output += "00000000"

	// reverse all bytes
	blen := 8
	if bytesep {
		blen++
	}
	nbytes := (len(output) - len(prefix) + 1) / blen

	for i := 0; i < nbytes; i++ {
		start := len(prefix) + i*blen
		end := start + 8
		bs := []byte(output[start:end])
		for i := 0; i < len(bs)/2; i++ {
			bs[i], bs[len(bs)-i-1] = bs[len(bs)-i-1], bs[i]
		}
		output = output[:start] + string(bs) + output[end:]
	}

	fmt.Println("output:", output)

	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer wait.Done()

		// read from server to stdout
		for {
			// byte by byte
			b, err := b.ReadByte()
			if err != nil {
				if err == io.EOF {
					fmt.Println("EOF")
					os.Exit(1)
				}
				panic(err)
			}
			bs := []byte{b}
			s := string(bs)
			fmt.Print(s)
		}
	}()

	// send output to server

	time.Sleep(100 * time.Millisecond)

	output += "\n"

	_, err = conn.Write([]byte(output))
	if err != nil {
		panic(err)
	}

	// wait for server to close connection
	fmt.Println("waiting")

	wait.Wait()

}
