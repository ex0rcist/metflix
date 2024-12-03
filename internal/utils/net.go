package utils

import (
	"fmt"
	"net"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
)

// Get outbound IP address of this host
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logging.LogError(err, "error closing utils.GetOutboundIP()")
		}
	}()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return net.IP{}, fmt.Errorf("utils.GetOutboundIP(): %w", entities.ErrUnexpected)
	}

	return localAddr.IP, nil
}
