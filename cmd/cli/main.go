package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	useragent := "vrc-social-notifier/0.1.0"

	subscriber := &streaming.VRChatStreamingSubscriber{}
	subscriber.OnFriendActive = func(event streaming.FriendActiveEvent) {
		fmt.Printf("FriendActive: %s\n", event.User.DisplayName)
	}
	subscriber.OnFriendOnline = func(event streaming.FriendOnlineEvent) {
		fmt.Printf("FriendOnline: %s\n", event.User.DisplayName)
	}
	subscriber.OnFriendOffline = func(event streaming.FriendOfflineEvent) {
		fmt.Printf("FriendOffline: %s\n", event.UserId)
	}

	if err := streaming.Subscribe(ctx, authToken, useragent, subscriber); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
