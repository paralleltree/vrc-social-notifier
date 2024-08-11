package feat

import (
	"context"

	"github.com/paralleltree/vrc-social-notifier/streaming"
)

type FriendStatusChangedEvent struct {
	PreviousStatus string
	CurrentStatus  string
	User           streaming.User
}

func NotifyFriendStatusChange(
	ctx context.Context,
	friendLocationCh <-chan streaming.FriendLocationEvent,
	friendUpdateCh <-chan streaming.FriendUpdateEvent,
) <-chan FriendStatusChangedEvent {
	notifyCh := make(chan FriendStatusChangedEvent)
	statusMap := map[string]string{}
	onStatusUpdate := func(user streaming.User) {
		if prevStatus, ok := statusMap[user.ID]; ok {
			if prevStatus != user.Status {
				notifyCh <- FriendStatusChangedEvent{
					PreviousStatus: prevStatus,
					CurrentStatus:  user.Status,
					User:           user,
				}
			}
		}
		statusMap[user.ID] = user.Status
	}

	go func() {
		defer close(notifyCh)

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

	return notifyCh
}
