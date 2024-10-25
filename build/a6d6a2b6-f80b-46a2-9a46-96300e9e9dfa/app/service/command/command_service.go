package command

import (
	"context"
	"noah/client/app/service"
	"os/exec"
	"runtime"
	"time"
)

type Service struct{}

func NewService() service.Command {
	return &Service{}
}

func (t Service) Run(command string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case `windows`:
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	case `linux`:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	case `darwin`:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	default:
		return nil, service.ErrUnsupportedPlatform
	}

	result, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return nil, service.ErrDeadlineExceeded
		}
		return result, nil
	}
	return result, nil
}
