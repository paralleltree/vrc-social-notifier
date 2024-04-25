package xsoverlay

import (
	"encoding/json"
	"fmt"
	"net"
)

func SendNotification(message Notification) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("localhost:%d", 42069))
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("dial udp: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(payload); err != nil {
		return fmt.Errorf("write payload: %w", err)
	}

	return nil
}
