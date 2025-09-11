package conn

import (
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sureshkrishnan-v/lazynet/internal/models"
)

// FetchConnections fetches all network connections and enriches them with process info.
func FetchConnections() ([]models.Connection, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	var result []models.Connection
	procCache := make(map[int32]*process.Process)

	for _, c := range conns {
		var conn models.Connection
		conn, procCache = models.FromNetConnectionStat(c, procCache)
		result = append(result, conn)
	}

	return result, nil
}

