package feat

import (
	"context"

	"github.com/paralleltree/vrc-social-notifier/streaming"
)

type FriendJoiningEvent struct {
	User streaming.User
}

func NotifyFriendJoining(
	ctx context.Context,
	userLocationCh <-chan streaming.UserLocationEvent,
	friendLocationCh <-chan streaming.FriendLocationEvent,
) <-chan FriendJoiningEvent {
	notifyCh := make(chan FriendJoiningEvent)
	currentLocation := ""

	go func() {
		defer close(notifyCh)
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
					notifyCh <- FriendJoiningEvent{
						User: friendLocation.User,
					}
				}
			}
		}
	}()

	return notifyCh
}
