package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math"

	"github.com/marcfyk/hackattic/internal/http_client"
)

const problem = "help_me_unpack"

type Input struct {
	Bytes string `json:"bytes"`
}

type Output struct {
	Int             int32   `json:"int"`
	UInt            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float32 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

func run() error {
	input, err := http_client.GetProblemInput[Input](problem)
	if err != nil {
		return err
	}
	dataBytes, err := base64.StdEncoding.DecodeString(input.Bytes)
	if err != nil {
		return err
	}

	result := Output{
		Int:             int32(binary.LittleEndian.Uint32(dataBytes[:4])),
		UInt:            binary.LittleEndian.Uint32(dataBytes[4:8]),
		Short:           int16(binary.LittleEndian.Uint16(dataBytes[8:10])),
		Float:           math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[12:16])),
		Double:          math.Float64frombits(binary.LittleEndian.Uint64(dataBytes[16:24])),
		BigEndianDouble: math.Float64frombits(binary.BigEndian.Uint64(dataBytes[24:])),
	}

	output, err := http_client.SubmitSolution(problem, result)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
