package feat

import (
	"context"
	"fmt"

	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/xsoverlay"
)

func NotifyFriendJoining(
	ctx context.Context,
	userLocationCh <-chan streaming.UserLocationEvent,
	friendLocationCh <-chan streaming.FriendLocationEvent,
	notifyCh chan<- xsoverlay.Notification,
) {
	currentLocation := ""
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
