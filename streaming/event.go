package streaming

type User struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
	InstanceID  string `json:"instanceId"`
}

type FriendActiveEvent struct {
	UserId string `json:"userId"`
	User   User   `json:"user"`
}

type FriendOnlineEvent struct {
	UserId string `json:"userId"`
	User   User   `json:"user"`
}

type FriendOfflineEvent struct {
	UserId string `json:"userId"`
}
