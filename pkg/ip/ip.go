package ip

import (
	"errors"
	"net"
)

var ErrDialingIp = errors.New("failed to get ip")

func GetLANIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.Join(ErrDialingIp, err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	localIP := localAddr.IP.String()

	return localIP, nil
}
