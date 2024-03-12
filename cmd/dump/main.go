package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paralleltree/vrc-social-notifier/streaming"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}

func run(ctx context.Context) error {
	authToken := os.Getenv("VRC_AUTH_TOKEN")
	if authToken == "" {
		return fmt.Errorf("VRC_AUTH_TOKEN is required")
	}
	useragent := "paltee.dev/0.1.0"

	subscriber := &streaming.VRChatStreamingSubscriber{}
	subscriber.OnMessageReceived = func(event string) {
		if err := processMessage(event); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", event)
			fmt.Fprintf(os.Stderr, "process message: %v\n", err)
		}
	}

	if err := streaming.Subscribe(ctx, authToken, useragent, subscriber); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func processMessage(msg string) error {
	payload := struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}{}
	if err := json.Unmarshal([]byte(msg), &payload); err != nil {
		return fmt.Errorf("unmarshal message: %w", err)
	}

	content := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload.Content), &content); err != nil {
		return fmt.Errorf("unmarshal content: %w", err)
	}
	contentMap := convertMapKeyToFlat("", content)

	logLine := map[string]interface{}{
		"time":         time.Now(),
		"raw":          msg,
		"message.type": payload.Type,
	}

	for k, v := range contentMap {
		logLine[fmt.Sprintf("message.content.%s", k)] = v
	}

	bytes, err := json.Marshal(logLine)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	fmt.Println(string(bytes))
	return nil
}

func convertMapKeyToFlat(parentKey string, v map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	var walkAndConvert func(parentKey string, v map[string]interface{})
	walkAndConvert = func(parentKey string, v map[string]interface{}) {
		prefix := fmt.Sprintf("%s.", parentKey)
		if parentKey == "" {
			prefix = ""
		}
		for k, v := range v {
			resultKey := fmt.Sprintf("%s%s", prefix, k)
			switch v := v.(type) {
			case map[string]interface{}:
				walkAndConvert(resultKey, v)
			default:
				result[resultKey] = v
			}
		}
	}

	walkAndConvert(parentKey, v)
	return result
}
