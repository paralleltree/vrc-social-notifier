package feat_test

import (
	"context"
	"strings"
	"testing"

	"github.com/paralleltree/vrc-social-notifier/feat"
	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/xsoverlay"
)

func TestNotifyFriendStatusChange(t *testing.T) {
	// arrange
	friendID := "usr_friend"
	friendName := "TestFriend_deadbeef"

	ctx := context.Background()
	friendLocationCh := make(chan streaming.FriendLocationEvent)
	friendUpdateCh := make(chan streaming.FriendUpdateEvent)
	notifyCh := make(chan xsoverlay.Notification)

	// act
	feat.NotifyFriendStatusChange(ctx, friendLocationCh, friendUpdateCh, notifyCh)

	go func() {
		friendLocationCh <- streaming.FriendLocationEvent{
			UserId: friendID,
			User: streaming.User{
				ID:          friendID,
				DisplayName: friendName,
				Status:      "busy",
			},
		}
		friendLocationCh <- streaming.FriendLocationEvent{
			UserId: friendID,
			User: streaming.User{
				ID:          friendID,
				DisplayName: friendName,
				Status:      "ask me",
			},
		}
	}()

	result := <-notifyCh

	// assert
	if !strings.Contains(result.Title, friendName) {
		t.Fatalf("unexpected notification: string `%s` was not found in notification title, got %s", friendName, result.Title)
	}
	if !strings.Contains(result.Title, "ask me") {
		t.Fatalf("unexpected notification: string `%s` was not found in notification title, got %s", "ask me", result.Title)
	}
}
