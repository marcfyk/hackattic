package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/marcfyk/hackattic/internal/http_client"
)

const problem = "reading_qr"

type Input struct {
	ImageURL string `json:"image_url"`
}

type Output struct {
	Code string `json:"code"`
}

func getQRCode(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	img, err := png.Decode(resp.Body)
	if err != nil {
		return "", err
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	result, err := qrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}
func run() error {
	input, err := http_client.GetProblemInput[Input](problem)
	if err != nil {
		return err
	}
	code, err := getQRCode(input.ImageURL)
	if err != nil {
		return err
	}
	resp, err := http_client.SubmitSolution(problem, Output{Code: code})
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
