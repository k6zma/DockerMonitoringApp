package domain

type PingResult struct {
	IP       string
	Name     string
	Status   string
	Success  bool
	PingTime int64
}

type ContainerInfo struct {
	IP     string
	Name   string
	Status string
}
