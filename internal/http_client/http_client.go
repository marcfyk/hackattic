package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
)

const (
	Domain         = "https://hackattic.com"
	EnvAccessToken = "HACKATTIC_ACCESS_TOKEN"
)

var ErrNoAccessTokenFound = fmt.Errorf(
	"no access token found on system, please set environment variable $%s",
	EnvAccessToken)

func GetProblemInput[A any](problem string) (*A, error) {
	accessToken, ok := os.LookupEnv(EnvAccessToken)
	if !ok {
		return nil, ErrNoAccessTokenFound
	}

	location, err := url.Parse(Domain)
	if err != nil {
		return nil, err
	}
	location.Path = fmt.Sprintf("/challenges/%s/problem", problem)

	queries := location.Query()
	queries.Add("access_token", accessToken)
	location.RawQuery = queries.Encode()

	resp, err := http.Get(location.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()

	var input A
	if err := decoder.Decode(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func SubmitSolution[A any](problem string, data A) ([]byte, error) {
	accessToken, ok := os.LookupEnv(EnvAccessToken)
	if !ok {
		return nil, ErrNoAccessTokenFound
	}

	location, err := url.Parse(Domain)
	if err != nil {
		return nil, err
	}
	location.Path = fmt.Sprintf("/challenges/%s/solve", problem)

	queries := location.Query()
	queries.Add("access_token", accessToken)
	location.RawQuery = queries.Encode()

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		location.String(),
		mime.TypeByExtension(".json"),
		bytes.NewReader(byteData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
