package sd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"net/http"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}

func DiskCheck(c *gin.Context) {
	usage, _ := disk.Usage("/")

	usedMB := int(usage.Used) / MB
	usedGB := int(usage.Used) / GB
	totalMB := int(usage.Total) / MB
	totalGB := int(usage.Total) / GB

	usedPercent := int(usage.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusOK
		text = "CRITICAL"
	} else if usedPercent > 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf(
		"%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%",
		text,
		usedMB,
		usedGB,
		totalMB,
		totalGB,
		usedPercent,
	)
	c.String(status, "\n"+message)
}

func CPUCheck(c *gin.Context) {
	cores, _ := cpu.Counts(false)

	avg, _ := load.Avg()

	l1 := avg.Load1
	l5 := avg.Load5
	l15 := avg.Load15

	status := http.StatusOK

	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s -Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)

	c.String(status, "\n"+message)
}

func RAMCheck(c *gin.Context) {
	usage, _ := mem.VirtualMemory()

	usedMB := int(usage.Used) / MB
	usedGB := int(usage.Used) / GB
	totalMB := int(usage.Total) / MB
	totalGB := int(usage.Total) / GB
	usedPercent := int(usage.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	c.String(status, "\n"+message)

}
