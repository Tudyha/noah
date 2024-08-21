package middleware

import (
	"github.com/robfig/cron/v3"
	"noah/internal/server/service"
)

func LoadCron() error {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		service.GetDeviceService().ScheduleUpdateStatus()
	})

	if err != nil {
		c.Stop()
		return err
	}

	c.Start()

	return nil
}
