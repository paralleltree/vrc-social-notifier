package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type VRChatStreamingSubscriber struct {
	OnConnected       func()
	OnMessageReceived func(string)
	OnError           func(string, error)
	OnUserLocation    func(UserLocationEvent)
	OnFriendActive    func(FriendActiveEvent)
	OnFriendOnline    func(FriendOnlineEvent)
	OnFriendOffline   func(FriendOfflineEvent)
	OnFriendLocation  func(FriendLocationEvent)
}

func Subscribe(ctx context.Context, authToken string, useragent string, subscriber *VRChatStreamingSubscriber) <-chan struct{} {
	ch := connectToVRChatStreaming(ctx, authToken, useragent, 5*time.Second)
	connClosed := make(chan struct{})
	onError := func(message string, err error) {
		if subscriber.OnError != nil {
			subscriber.OnError(message, err)
		}
	}
	go func() {
		defer close(connClosed)
		for msg := range ch {
			if msg.Err != nil {
				if subscriber.OnError != nil {
					onError(msg.Value, msg.Err)
				}
				continue
			}
			if msg.Connected {
				if subscriber.OnConnected != nil {
					subscriber.OnConnected()
				}
				continue
			}
			if err := subscriber.processVRChatEvent(msg.Value); err != nil {
				onError(msg.Value, err)
			}
		}
	}()
	return connClosed
}

func (s *VRChatStreamingSubscriber) processVRChatEvent(msg string) error {
	if s.OnMessageReceived != nil {
		s.OnMessageReceived(msg)
	}

	meta := struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}{}
	if err := json.Unmarshal([]byte(msg), &meta); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	switch meta.Type {
	case "user-location":
		payload := UserLocationEvent{}
		if err := json.Unmarshal([]byte(meta.Content), &payload); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		if s.OnUserLocation != nil {
			s.OnUserLocation(payload)
		}

	case "friend-active":
		payload := FriendActiveEvent{}
		if err := json.Unmarshal([]byte(meta.Content), &payload); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		if s.OnFriendActive != nil {
			s.OnFriendActive(payload)
		}

	case "friend-online":
		payload := FriendOnlineEvent{}
		if err := json.Unmarshal([]byte(meta.Content), &payload); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		if s.OnFriendOnline != nil {
			s.OnFriendOnline(payload)
		}

	case "friend-offline":
		payload := FriendOfflineEvent{}
		if err := json.Unmarshal([]byte(meta.Content), &payload); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		if s.OnFriendOffline != nil {
			s.OnFriendOffline(payload)
		}

	case "friend-location":
		payload := FriendLocationEvent{}
		if err := json.Unmarshal([]byte(meta.Content), &payload); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
		if s.OnFriendLocation != nil {
			s.OnFriendLocation(payload)
		}

	default:
	}

	return nil
}
