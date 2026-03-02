package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func Get(url string) (map[string]any, error) {
	var data map[string]any
	httpClient := http.Client{}
	resp, err := httpClient.Get(url)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	json.Unmarshal(body, &data)

	if err != nil {
		return data, err
	}

	return data, nil
}

func GetRaw(url string) ([]byte, error) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}
