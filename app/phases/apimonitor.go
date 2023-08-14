package pkg

import "time"

type ApiMonitor struct {
	Endpoint string
}

func (l ApiMonitor) WaitForClusterReady(timeout time.Duration) {
	return
}

func (l ApiMonitor) WaitForMastersReady(timeout time.Duration){
}

func (l ApiMonitor) WaitForWorkersReady(timeout time.Duration){
	return
}

