package task

import "net/http"

type Dashboard struct {
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (d *Dashboard) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
