package middleware

import (
	"github.com/robfig/cron/v3"
	"noah/internal/server/service"
)

func LoadCron() error {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		service.GetClientService().ScheduleUpdateStatus()
	})

	//每天0点执行一次
	_, err = c.AddFunc("0 0 * * *", func() {
		service.GetClientService().CleanSystemInfo()
	})

	if err != nil {
		c.Stop()
		return err
	}

	c.Start()

	return nil
}
