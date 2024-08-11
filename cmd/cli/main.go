package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/paralleltree/vrc-social-notifier/feat"
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

// Returns function that creates new channel that receives elements from source channel
func makeChGenerator[T any](source <-chan T) func() chan T {
	chs := []chan T{}

	go func() {
		defer func() {
			// close all created channels when source channel is closed
			for _, ch := range chs {
				close(ch)
			}
		}()

		for e := range source {
			for _, ch := range chs {
				ch <- e
			}
		}
	}()

	return func() chan T {
		ch := make(chan T)
		chs = append(chs, ch)
		return ch
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
	friendUpdateCh := make(chan streaming.FriendUpdateEvent)

	makeSendUserLocationCh := makeChGenerator(userLocationCh)
	makeSendFriendLocationCh := makeChGenerator(friendLocationCh)
	makeSendFriendUpdateCh := makeChGenerator(friendUpdateCh)

	go func() {
		friendJoiningCh := feat.NotifyFriendJoining(ctx, makeSendUserLocationCh(), makeSendFriendLocationCh())
		for v := range friendJoiningCh {
			n := xsoverlay.NewNotificationBuilder().
				SetTitle(fmt.Sprintf("Friend Joining: %s", v.User.DisplayName)).
				Build()
			inboxCh <- n
		}
	}()
	go func() {
		friendStatusChangedCh := feat.NotifyFriendStatusChange(ctx, makeSendFriendLocationCh(), makeSendFriendUpdateCh())
		for v := range friendStatusChangedCh {
			n := xsoverlay.NewNotificationBuilder().
				SetTitle(fmt.Sprintf("%s's status changed to %s", v.User.DisplayName, v.User.Status)).
				Build()
			inboxCh <- n
		}
	}()

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
	subscriber.OnFriendUpdate = func(event streaming.FriendUpdateEvent) {
		fmt.Printf("FriendUpdate: %s(%s: %s)\n", event.User.DisplayName, event.User.Status, event.User.StatusDescription)
		friendUpdateCh <- event
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
