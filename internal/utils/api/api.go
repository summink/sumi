package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

var downloadClient = &http.Client{
	Timeout: 10 * time.Minute,
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

type ProgressCallback func(downloaded, total int64)

func GetWithProgress(url string, onProgress ProgressCallback) ([]byte, error) {
	const maxRetries = 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		data, err := downloadWithProgress(url, onProgress)
		if err == nil {
			return data, nil
		}
		lastErr = err
		if attempt < maxRetries {
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}

	return nil, lastErr
}

func downloadWithProgress(url string, onProgress ProgressCallback) ([]byte, error) {
	resp, err := downloadClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	total := resp.ContentLength

	var data []byte
	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
			downloaded += int64(n)
			if onProgress != nil {
				onProgress(downloaded, total)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func CalculateSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
