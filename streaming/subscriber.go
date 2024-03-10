package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type VRChatStreamingSubscriber struct {
	OnFriendActive  func(FriendActiveEvent)
	OnFriendOnline  func(FriendOnlineEvent)
	OnFriendOffline func(FriendOfflineEvent)
}

func Subscribe(ctx context.Context, authToken string, useragent string, subscriber *VRChatStreamingSubscriber) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := connectToVRChatStreaming(ctx, authToken, useragent, 5*time.Second)
	for msg := range ch {
		if msg.Err != nil {
			return msg.Err
		}
		if err := subscriber.processVRChatEvent(msg.Value); err != nil {
			return err
		}
	}
	return nil
}

func (s *VRChatStreamingSubscriber) processVRChatEvent(msg string) error {
	meta := struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}{}
	if err := json.Unmarshal([]byte(msg), &meta); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	switch meta.Type {
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

	default:
	}

	return nil
}
