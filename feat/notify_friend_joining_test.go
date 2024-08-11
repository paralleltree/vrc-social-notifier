package feat_test

import (
	"context"
	"testing"

	"github.com/paralleltree/vrc-social-notifier/feat"
	"github.com/paralleltree/vrc-social-notifier/streaming"
)

func TestNotifyFriendJoining(t *testing.T) {
	// arrange
	userID := "usr_my"
	friendID := "usr_friend"
	friendName := "TestFriend_deadbeef"
	locationID := "loc_test"

	wantResult := feat.FriendJoiningEvent{
		User: streaming.User{
			ID:          friendID,
			DisplayName: friendName,
		},
	}
	wantResults := []feat.FriendJoiningEvent{wantResult}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	friendLocationCh := make(chan streaming.FriendLocationEvent)
	userLocationCh := make(chan streaming.UserLocationEvent)

	// act
	notifyCh := feat.NotifyFriendJoining(ctx, userLocationCh, friendLocationCh)

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
				ID:          friendID,
				DisplayName: friendName,
			},
		}
	}()

	gotResults := consumeChannel(notifyCh)

	// assert
	if len(wantResults) != len(gotResults) {
		t.Fatalf("unexpected notification size: want %d, but got %d", len(wantResults), len(gotResults))
	}

	for i, gotResult := range gotResults {
		if wantResults[i].User.ID != gotResult.User.ID {
			t.Fatalf("unexpected user: want %s, but got %s", wantResults[i].User.ID, gotResult.User.ID)
		}
	}
}
