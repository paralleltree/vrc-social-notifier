package streaming

import (
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/websocket"
)

type streamingEvent struct {
	Value string
	Err   error
}

func connectToWebSocket(url string, useragent string) (*websocket.Conn, func(), error) {
	conf, err := websocket.NewConfig(url, "http://localhost/")
	conf.Dialer = &net.Dialer{}
	if err != nil {
		return nil, nil, fmt.Errorf("create websocket config: %w", err)
	}
	conf.Header.Add("User-Agent", useragent)
	ws, err := websocket.DialConfig(conf)
	if err != nil {
		return nil, nil, fmt.Errorf("dial websocket: %w", err)
	}
	cancel := func() {
		ws.Close()
	}
	return ws, cancel, nil
}

func connectToVRChatStreaming(ctx context.Context, authToken string, useragent string, reconnectInterval time.Duration) <-chan streamingEvent {
	ch := make(chan streamingEvent)
	url := fmt.Sprintf("wss://pipeline.vrchat.cloud:443/?authToken=%s", authToken)
	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				ws, cancel, err := connectToWebSocket(url, useragent)
				if err != nil {
					ch <- streamingEvent{Err: err}

					// wait for reconnect
					select {
					case <-ctx.Done():
						return
					case <-time.After(reconnectInterval):
						continue
					}
				}
				defer cancel()

				for {
					select {
					case <-ctx.Done():
						return

					default:
						var msg string
						if err := websocket.Message.Receive(ws, &msg); err != nil {
							ch <- streamingEvent{Err: err}
							break
						}
						ch <- streamingEvent{Value: msg}
					}
				}
			}
		}
	}()
	return ch
}
