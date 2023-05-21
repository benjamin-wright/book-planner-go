package nats

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

func Run(name string, port int64) func() {
	os.Setenv("EVENTS_URL", "nats://127.0.0.1:4222")

	// Try to remove any existing containers
	exec.Command("docker", "stop", name).Run()

	image := "nats"

	cmdString := fmt.Sprintf(
		"run --rm -d -p %d:%d --name %s %s",
		port, port, name, image,
	)

	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	if err := cmd.Run(); err != nil {
		zap.S().Fatalf("Failed to start nats container: %+v", err)
	}

	return func() {
		cmd := exec.Command("docker", "stop", name)
		if err := cmd.Run(); err != nil {
			zap.S().Fatalf("Failed to stop nats container: %+v", err)
		}
	}
}
