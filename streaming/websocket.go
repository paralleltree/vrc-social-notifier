package streaming

import (
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/websocket"
)

type streamingEvent struct {
	Connected bool
	Value     string
	Err       error
}

func connectToWebSocket(url string, useragent string) (*websocket.Conn, error) {
	conf, err := websocket.NewConfig(url, "http://localhost/")
	conf.Dialer = &net.Dialer{}
	if err != nil {
		return nil, fmt.Errorf("create websocket config: %w", err)
	}
	conf.Header.Add("User-Agent", useragent)
	ws, err := websocket.DialConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("dial websocket: %w", err)
	}
	return ws, nil
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
				ws, err := connectToWebSocket(url, useragent)
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

				ch <- streamingEvent{Connected: true}
				connClosed := make(chan struct{})
				go func() {
					defer close(connClosed)
					defer ws.Close()

					for {
						var msg string
						if err := websocket.Message.Receive(ws, &msg); err != nil {
							ch <- streamingEvent{Err: err}
							return
						}
						ch <- streamingEvent{Value: msg}
					}
				}()

				select {
				case <-ctx.Done():
					ws.Close()

				case <-connClosed:
					continue
				}
			}
		}
	}()
	return ch
}
