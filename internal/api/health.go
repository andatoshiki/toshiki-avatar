package api

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"
)

var startTime = time.Now()

type HealthResponse struct {
	Uptime string `json:"uptime"`
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	result := ""
	if days > 0 {
		result += plural(days, "day") + " "
	}
	if hours > 0 {
		result += plural(hours, "hour") + " "
	}
	if minutes > 0 {
		result += plural(minutes, "min") + " "
	}
	result += plural(seconds, "sec")
	return result
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	resp := HealthResponse{Uptime: formatUptime(uptime)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
