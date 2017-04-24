package common

import (
	"fmt"
	"logger"
	"net"
	"os"
	"os/exec"
	"time"
)

var (
	Logger logger.Logger
)

type UsageTracker interface {
	StartTrackingUsage()
	MeasureUsage() (usage float64, measurementDuration time.Duration)
}

func LocalIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func LogCmd(cmd string) {
	Logger.Debug(cmd)
}

func ExecShell(cmd string) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		LogCmd(fmt.Sprintf("Running System Command: %s, Status:Error,  Detail is: %v \n", cmd, err))
		return "", err
	}
	result := string(out[:len(out)])
	LogCmd(fmt.Sprintf("Running System Command: %s, Status:Correct, Output: %s \n", cmd, result))
	return result, nil
}

func GenerateUUID() string {
	file, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	file.Read(b)
	file.Close()

	uuid := fmt.Sprintf("%x", b)
	return uuid
}
