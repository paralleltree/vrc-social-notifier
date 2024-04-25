package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type betterStackHeartbeat struct {
	endpoint string
}

func NewBetterStackHeartbeat(endpoint string) *betterStackHeartbeat {
	return &betterStackHeartbeat{endpoint: endpoint}
}

func (b *betterStackHeartbeat) SendHeartbeat(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", b.endpoint, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	return nil
}
