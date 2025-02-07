package domain

type PingResult struct {
	IP       string
	Success  bool
	PingTime float64
}
