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
	friendUpdateCh <-chan streaming.FriendUpdateEvent,
	notifyCh chan<- xsoverlay.Notification,
) {
	statusMap := map[string]string{}
	onStatusUpdate := func(user streaming.User) {
		if prevStatus, ok := statusMap[user.ID]; ok {
			if prevStatus != user.Status {
				n := xsoverlay.NewNotificationBuilder().
					SetTitle(fmt.Sprintf("%s's status changed to %s", user.DisplayName, user.Status)).
					Build()
				notifyCh <- n
			}
		}
		statusMap[user.ID] = user.Status
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case friendLocation := <-friendLocationCh:
				onStatusUpdate(friendLocation.User)

			case friendUpdate := <-friendUpdateCh:
				onStatusUpdate(friendUpdate.User)
			}
		}
	}()
}
