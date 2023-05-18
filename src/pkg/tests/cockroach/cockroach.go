package cockroach

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

func Run(name string, port int64) func() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	os.Setenv("POSTGRES_HOST", "127.0.0.1")
	os.Setenv("POSTGRES_PORT", strconv.FormatInt(port, 10))
	os.Setenv("POSTGRES_USER", "root")
	os.Setenv("POSTGRES_NAME", "defaultdb")

	// Try to remove any existing containers
	exec.Command("docker", "stop", name).Run()

	image := "cockroachdb/cockroach:v22.2.8"
	args := "--logtostderr start-single-node --insecure --listen-addr 0.0.0.0:" + strconv.FormatInt(port, 10)

	cmdString := fmt.Sprintf(
		"run --rm -d -p %d:%d --name %s %s %s",
		port, port, name, image, args,
	)

	cmd := exec.Command("docker", strings.Split(cmdString, " ")...)
	if err := cmd.Run(); err != nil {
		zap.S().Fatalf("Failed to start cockroach container: %+v", err)
	}

	return func() {
		cmd := exec.Command("docker", "stop", name)
		if err := cmd.Run(); err != nil {
			zap.S().Fatalf("Failed to stop cockroach container: %+v", err)
		}
	}
}

func Migrate(path string) {
	cfg, err := postgres.ConfigFromEnv()
	if err != nil {
		zap.S().Fatalf("Failed to get connection details: %+v", err)
	}

	conn, err := postgres.Connect(cfg)
	if err != nil {
		zap.S().Fatalf("Failed to connect to cockroach: %+v", err)
	}
	defer conn.Close(context.TODO())

	data, err := os.ReadFile(path)
	if err != nil {
		zap.S().Fatalf("Failed to read migration file: %+v", err)
	}

	_, err = conn.Exec(context.TODO(), string(data))
	if err != nil {
		zap.S().Fatalf("Failed to run migration: %+v", err)
	}

	zap.S().Infof("Ran migration: %s", path)
}
