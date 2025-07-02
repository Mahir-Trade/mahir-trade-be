package controller

import (
	"context"
	"log/slog"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/internal/app/service"
	"sync"

	"go.uber.org/dig"
)

type (
	SchedulerCtrl interface{}

	SchedulerCtrlImpl struct {
		dig.In

		Cron    infra.Schedule
		UserSvc service.UserSvc
	}
)

func NewSchedulerCtrl(impl SchedulerCtrlImpl) SchedulerCtrl {
	wg := sync.WaitGroup{}
	jobs := []infra.RegisterJobs{
		{
			Spec: "*/10 * * * *",
			Cmd: func(ctx context.Context) {
				err := impl.UserSvc.UpdateMembership(ctx)
				if err != nil {
					slog.ErrorContext(ctx, "scheduler", "failed to update user status", err)
				}
			},
		},
	}

	impl.Cron.Run(&wg, jobs)
	return &impl
}
