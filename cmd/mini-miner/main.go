package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/marcfyk/hackattic/internal/http_client"
)

const (
	problem = "mini_miner"
)

type Input struct {
	Difficulty int   `json:"difficulty"`
	Block      Block `json:"block"`
}

func (i Input) String() string {
	return fmt.Sprintf(`{"block":%s,"difficulty":%d}`, i.Block, i.Difficulty)
}

type Block struct {
	Data  [][]interface{} `json:"data"`
	Nonce *int            `json:"nonce"`
}

func (b Block) hashSHA256() ([32]byte, error) {
	dataBytes, err := json.Marshal(b)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(dataBytes), nil
}

func (b *Block) computeNonce(leadingZeroes int) (int, error) {
	nonce := 0
	b.Nonce = &nonce
	for {
		hash, err := b.hashSHA256()
		if err != nil {
			return 0, err
		}
		if hasLeadingZeroBits(hash[:], leadingZeroes) {
			return *b.Nonce, nil
		}
		*b.Nonce++
	}
}

func (b Block) String() string {
	nonce := "nil"
	if b.Nonce != nil {
		nonce = strconv.Itoa(*b.Nonce)
	}
	return fmt.Sprintf(`{"data":%s,"nonce":%s}`, b.Data, nonce)
}

func hasLeadingZeroBits(hash []byte, count int) bool {
	if len(hash)*8 < count {
		return false
	}

	index := 0
	for count >= 8 {
		if hash[index] != 0 {
			return false
		}
		index++
		count -= 8
	}
	if count > 0 {
		shifted := hash[index] >> (8 - count)
		return shifted == 0
	}
	return true
}

type Output struct {
	Nonce int `json:"nonce"`
}

func run() error {
	input, err := http_client.GetProblemInput[Input](problem)
	if err != nil {
		return err
	}
	nonce, err := input.Block.computeNonce(input.Difficulty)
	if err != nil {
		return err
	}
	resp, err := http_client.SubmitSolution(problem, Output{Nonce: nonce})
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}
