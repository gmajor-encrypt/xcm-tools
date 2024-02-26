package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func HttpPost(ctx context.Context, data []byte, endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func HttpGet(ctx context.Context, endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}
