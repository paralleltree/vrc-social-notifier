package streaming

type FriendOnlineEvent struct {
	UserId string `json:"userId"`
	User   struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		InstanceID  string `json:"instanceId"`
	} `json:"user"`
}

type FriendOfflineEvent struct {
	UserId string `json:"userId"`
}
