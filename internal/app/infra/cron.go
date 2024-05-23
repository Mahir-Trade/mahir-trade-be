package infra

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/dig"
)

type (
	ScheduleImpl struct {
		dig.Out

		CronJob *cron.Cron
	}

	RegisterJobs struct {
		Spec string
		Cmd  func(ctx context.Context)
	}

	Schedule interface {
		Run(wg *sync.WaitGroup, jobs []RegisterJobs)
		Stop()
	}
)

func NewScheduler() Schedule {
	return &ScheduleImpl{
		CronJob: initCron(),
	}
}

func (s *ScheduleImpl) Run(wg *sync.WaitGroup, jobs []RegisterJobs) {
	jobContext := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigs
		slog.Info("SIGTERM", "received signal", sig)
		s.CronJob.Stop()
		slog.InfoContext(jobContext, "SCHEDULER", "scheduler stopped at", time.Now())
	}()

	for _, job := range jobs {
		go func(spec string, cmd func(ctx context.Context)) {
			s.addFunc(spec, func() {
				cmd(jobContext)
			})
		}(job.Spec, job.Cmd)
	}
	s.CronJob.Start()
	slog.InfoContext(jobContext, "SCHEDULER", "scheduler started at", time.Now())
}

func initCron() *cron.Cron {
	return cron.New(cron.WithLocation(time.Local))
}

func (s *ScheduleImpl) Stop() {
	s.CronJob.Stop()
}

func (s *ScheduleImpl) addFunc(spec string, cmd func()) {
	id, err := s.CronJob.AddFunc(spec, cmd)
	if err != nil {
		slog.Error("SCHEDULER", "failed to add job", err)
		panic("failed to add job")
	}
	slog.Info("JOB", "job added", id)
}
