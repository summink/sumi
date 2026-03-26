package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func Get(url string) (map[string]any, error) {
	var data map[string]any

	resp, err := httpClient.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func GetRaw(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
