package domain

type PingResult struct {
	IP       string `json:"ip_address"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Success  bool   `json:"success"`
	PingTime int64  `json:"ping_time"`
	LastPing string `json:"last_successful_ping"`
}

type ContainerInfo struct {
	IP     string
	Name   string
	Status string
}
