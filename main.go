package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	data := os.Args[1]

	err := json.Unmarshal([]byte(data), &in)
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

	// b1
	fmt.Println("00010001")
	// b2
	fmt.Println("00000110")
	// b3 / a1 MP
	fmt.Printf("%016b\n", in.ReportingPoint)
	// a2 PR 00000000 // 11?? dir /
	fmt.Printf("%02b\n", in.Priority)
	fmt.Printf("%02b\n", in.DirectionRequest)
	fmt.Printf("%04b\n", in.Line>>8)
	// a3 LN  line number
	fmt.Printf("%08b\n", in.Line&0xFF)
	// a4 KN run_number
	fmt.Printf("%07b\n", in.RunNumber)
	// a5 ZN dest
	fmt.Printf("%08b\n", in.DestinationNumber>>4)
	// a6 ZN / ZL
	fmt.Printf("%04b\n", in.DestinationNumber&0xF)
	fmt.Println("0")
	fmt.Printf("%03b\n", in.TrainLength)

}
