package command

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"noah/client/app/entitie"
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

func (t Service) GetProcessList() ([]entitie.Process, error) {
	//获取系统进程
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	//
	var result []entitie.Process
	for _, proc := range processes {
		if proc.Pid == 0 {
			fmt.Println("pid = 0 进程不显示")
			continue
		}
		name, _ := proc.Name()
		uid, _ := proc.Uids()
		gid, _ := proc.Gids()
		cmdline, _ := proc.Cmdline()
		username, _ := proc.Username()
		cpu, _ := proc.Percent(0)
		m, _ := proc.MemoryInfo()
		m_rss := uint64(0)
		if m != nil {
			m_rss = m.RSS
		}
		createTime, _ := proc.CreateTime()

		result = append(result, entitie.Process{
			Pid:        proc.Pid,
			Name:       name,
			Username:   username,
			Uids:       uid,
			Gids:       gid,
			Command:    cmdline,
			Cpu:        cpu,
			Memory:     m_rss,
			CreateTime: createTime,
		})
	}
	return result, nil
}

func (t Service) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}
