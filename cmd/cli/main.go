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
	subscriber.OnError = func(message string, err error) {
		fmt.Fprintf(os.Stderr, "%s\n", message)
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	subscriber.OnUserLocation = func(event streaming.UserLocationEvent) {
		fmt.Printf("UserLocation: user: %s, location: %s, instance: %s, world: %s\n", event.User.DisplayName, event.Location, event.Instance, event.WorldId)
	}
	subscriber.OnFriendActive = func(event streaming.FriendActiveEvent) {
		fmt.Printf("FriendActive: %s\n", event.User.DisplayName)
	}
	subscriber.OnFriendOnline = func(event streaming.FriendOnlineEvent) {
		fmt.Printf("FriendOnline: %s\n", event.User.DisplayName)
	}
	subscriber.OnFriendOffline = func(event streaming.FriendOfflineEvent) {
		fmt.Printf("FriendOffline: %s\n", event.UserId)
	}
	subscriber.OnFriendLocation = func(event streaming.FriendLocationEvent) {
		fmt.Printf("FriendLocation: user: %s, location: %s, travelingToLocation: %s, instance: %s\n", event.User.DisplayName, event.Location, event.TravelingToLocation, event.Instance)
	}

	connClosed := streaming.Subscribe(ctx, authToken, useragent, subscriber)
	<-ctx.Done()
	<-connClosed
	return nil
}
