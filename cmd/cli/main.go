package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/xsoverlay"
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

	inboxCh := make(chan xsoverlay.Notification)
	go func() {
		for message := range inboxCh {
			if err := xsoverlay.SendNotification(message); err != nil {
				fmt.Fprintf(os.Stderr, "xsoverlay error: %v\n", err)
			}
		}
	}()

	userLocationCh := make(chan streaming.UserLocationEvent)
	friendLocationCh := make(chan streaming.FriendLocationEvent)

	notifyFriendJoining(ctx, userLocationCh, friendLocationCh, inboxCh)

	subscriber := &streaming.VRChatStreamingSubscriber{}
	subscriber.OnError = func(message string, err error) {
		fmt.Fprintf(os.Stderr, "%s\n", message)
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	subscriber.OnUserLocation = func(event streaming.UserLocationEvent) {
		fmt.Printf("UserLocation: user: %s, location: %s, instance: %s, world: %s\n", event.User.DisplayName, event.Location, event.Instance, event.WorldId)
		userLocationCh <- event
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
		friendLocationCh <- event
	}

	connClosed := streaming.Subscribe(ctx, authToken, useragent, subscriber)
	<-ctx.Done()
	<-connClosed
	return nil
}

func notifyFriendJoining(
	ctx context.Context,
	userLocationCh <-chan streaming.UserLocationEvent,
	friendLocationCh <-chan streaming.FriendLocationEvent,
	notifyCh chan<- xsoverlay.Notification,
) {
	currentLocation := ""
	notifyCh <- xsoverlay.NewNotificationBuilder().SetTitle("VRC Social Notification enabled").Build()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case userLocation := <-userLocationCh:
				if userLocation.Location != currentLocation {
					currentLocation = userLocation.Location
				}

			case friendLocation := <-friendLocationCh:
				if friendLocation.Location == "traveling" && currentLocation == friendLocation.TravelingToLocation {
					n := xsoverlay.NewNotificationBuilder().
						SetTitle(fmt.Sprintf("Friend Joining: %s", friendLocation.User.DisplayName)).
						Build()
					notifyCh <- n
				}
			}
		}
	}()
}
