package utils

import (
	"github.com/shirou/gopsutil/net"
)

type ActivePort struct {
	Port int
	Pid  int
}

func GetRunningPorts() ([]*ActivePort, error) {
	connections, err := net.Connections("inet")
	if err != nil {
		return nil, err
	}

	out := make([]*ActivePort, 0)
	for _, cs := range connections {
		if cs.Status == "LISTEN" {
			out = append(out, &ActivePort{
				Port: int(cs.Laddr.Port),
				Pid:  int(cs.Pid),
			})
		}
	}

	return out, nil
}
