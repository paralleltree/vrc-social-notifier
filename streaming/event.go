package streaming

type User struct {
	DisplayName       string `json:"displayName"`
	ID                string `json:"id"`
	InstanceID        string `json:"instanceId"`
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type UserLocationEvent struct {
	UserId   string `json:"userId"`
	User     User   `json:"user"`
	Location string `json:"location"`
	Instance string `json:"instance"`
	WorldId  string `json:"worldId"`
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

type FriendLocationEvent struct {
	UserId string `json:"userId"`
	User   User   `json:"user"`

	Location string `json:"location"`

	TravelingToLocation string `json:"travelingToLocation"`

	Instance string `json:"instance"`

	WorldId string `json:"worldId"`

	CanRequestInvite bool `json:"canRequestInvite"`
}
