package feat_test

// t.Fatalf("unexpected notification: string `%s` was not found in notification title, got %s", friendName, result.Title)

import (
	"context"
	"testing"

	"github.com/paralleltree/vrc-social-notifier/feat"
	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/testlib"
)

func TestNotifyFriendStatusChange(t *testing.T) {
	// arrange
	friendID := "usr_friend"
	friendName := "TestFriend_deadbeef"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	friendLocationCh := make(chan streaming.FriendLocationEvent)
	friendUpdateCh := make(chan streaming.FriendUpdateEvent)

	wantResults := []feat.FriendStatusChangedEvent{
		{
			PreviousStatus: "busy",
			CurrentStatus:  "ask me",
			User: streaming.User{
				ID:          friendID,
				DisplayName: friendName,
			},
		},
	}

	// act
	notifyCh := feat.NotifyFriendStatusChange(ctx, friendLocationCh, friendUpdateCh)

	go func() {
		defer cancel()

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

	gotResults := testlib.ConsumeChannel(notifyCh)

	// assert
	if len(wantResults) != len(gotResults) {
		t.Fatalf("unexpected notification size: want %d, but got %d", len(wantResults), len(gotResults))
	}

	for i, gotResult := range gotResults {
		if wantResults[i].PreviousStatus != gotResult.PreviousStatus {
			t.Errorf("unexpected previous status: want %s, but got %s", wantResults[i].PreviousStatus, gotResult.PreviousStatus)
		}
		if wantResults[i].CurrentStatus != gotResult.CurrentStatus {
			t.Errorf("unexpected current status: want %s, but got %s", wantResults[i].CurrentStatus, gotResult.CurrentStatus)
		}
		if wantResults[i].User.ID != gotResult.User.ID {
			t.Errorf("unexpected user id: want %s, but got %s", wantResults[i].User.ID, gotResult.User.ID)
		}
		if wantResults[i].User.DisplayName != gotResult.User.DisplayName {
			t.Errorf("unexpected user name: want %s, but got %s", wantResults[i].User.DisplayName, gotResult.User.DisplayName)
		}
	}
}
