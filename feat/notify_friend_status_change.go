package feat

import (
	"context"
	"fmt"

	"github.com/paralleltree/vrc-social-notifier/streaming"
	"github.com/paralleltree/vrc-social-notifier/xsoverlay"
)

func NotifyFriendStatusChange(
	ctx context.Context,
	friendLocationCh <-chan streaming.FriendLocationEvent,
	notifyCh chan<- xsoverlay.Notification,
) {
	statusMap := map[string]string{}
	notifyCh <- xsoverlay.NewNotificationBuilder().SetTitle("NotifyFriendStatusChange enabled").Build()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case friendLocation := <-friendLocationCh:
				if prevStatus, ok := statusMap[friendLocation.User.ID]; ok {
					if prevStatus != friendLocation.User.Status {
						n := xsoverlay.NewNotificationBuilder().
							SetTitle(fmt.Sprintf("User %s status changed to %s", friendLocation.User.DisplayName, friendLocation.User.Status)).
							SetBody(friendLocation.User.StatusDescription).
							Build()
						notifyCh <- n
					}
				}
				statusMap[friendLocation.User.ID] = friendLocation.User.Status
			}
		}
	}()
}
