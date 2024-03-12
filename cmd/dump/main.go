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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
	subscriber.OnError = func(message string, err error) {
		fmt.Fprintf(os.Stderr, "%s\n", message)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		cancel()
	}

	connClosed := streaming.Subscribe(ctx, authToken, useragent, subscriber)
	<-ctx.Done()
	fmt.Fprintf(os.Stderr, "shutting down...\n")
	<-connClosed
	fmt.Fprintf(os.Stderr, "exiting\n")
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

	contentMap, err := convertMessageContentToFlatMap(payload.Type, payload.Content)
	if err != nil {
		return err
	}

	logLine := map[string]interface{}{
		"time":         time.Now(),
		"raw":          msg,
		"message.type": payload.Type,
	}

	for k, v := range contentMap {
		logLine[fmt.Sprintf("message.%s", k)] = v
	}

	bytes, err := json.Marshal(logLine)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	fmt.Println(string(bytes))
	return nil
}

func convertMessageContentToFlatMap(messageType, content string) (map[string]interface{}, error) {
	switch messageType {
	case "see-notification", "hide-notification":
		return map[string]interface{}{
			"content": content,
		}, nil

	default:
		contentMap := map[string]interface{}{}
		if err := json.Unmarshal([]byte(content), &contentMap); err != nil {
			return nil, fmt.Errorf("unmarshal content: %w", err)
		}
		return convertMapKeyToFlat("content", contentMap), nil
	}
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
