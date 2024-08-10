package feat_test

import (
	"context"
	"strings"
	"testing"

	"github.com/paralleltree/vrc-social-notifier/feat"
	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/xsoverlay"
)

func TestNotifyFriendJoining(t *testing.T) {
	// arrange
	userID := "usr_my"
	friendID := "usr_friend"
	friendName := "TestFriend_deadbeef"
	locationID := "loc_test"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	friendLocationCh := make(chan streaming.FriendLocationEvent)
	userLocationCh := make(chan streaming.UserLocationEvent)
	notifyCh := make(chan xsoverlay.Notification)

	// act
	feat.NotifyFriendJoining(ctx, userLocationCh, friendLocationCh, notifyCh)

	go func() {
		defer cancel()

		userLocationCh <- streaming.UserLocationEvent{
			UserId:   userID,
			Location: locationID,
		}
		friendLocationCh <- streaming.FriendLocationEvent{
			UserId:              friendID,
			Location:            "traveling",
			TravelingToLocation: locationID,
			User: streaming.User{
				DisplayName: friendName,
			},
		}
	}()

	result := <-notifyCh

	if !strings.Contains(result.Title, friendName) {
		t.Fatalf("unexpected notification: string `%s` was not found in notification title, got %s", friendName, result.Title)
	}
}
