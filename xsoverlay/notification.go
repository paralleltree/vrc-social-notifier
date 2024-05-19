package xsoverlay

type Notification struct {
	Type          int     `json:"messageType"`
	Timeout       float32 `json:"timeout"`   // Defalut 0.5
	Height        float32 `json:"height"`    // Default 175
	Opacity       float32 `json:"opacity"`   // Default 1
	Volume        float32 `json:"volume"`    // Default 0.7
	AudioPath     string  `json:"audioPath"` // path, default, error, waring
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	UseBase64Icon bool    `json:"useBase64Icon"`
	Icon          string  `json:"icon"` // base64 encoded image or path
	SourceApp     string  `json:"sourceApp"`
}

type notificationBuilder struct {
	Notification Notification
}

func NewNotificationBuilder() *notificationBuilder {
	return &notificationBuilder{
		Notification: Notification{
			Type:      1,
			Timeout:   2,
			Height:    175,
			Opacity:   1,
			Volume:    0.7,
			AudioPath: "default",
		},
	}
}

func (b *notificationBuilder) Build() Notification {
	return b.Notification
}

func (b *notificationBuilder) SetTitle(title string) *notificationBuilder {
	b.Notification.Title = title
	return b
}

func (b *notificationBuilder) SetBody(body string) *notificationBuilder {
	b.Notification.Content = body
	return b
}

func (b *notificationBuilder) SetContent(content string) *notificationBuilder {
	b.Notification.Content = content
	return b
}
