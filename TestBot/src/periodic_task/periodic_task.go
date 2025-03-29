package periodictask

import (
	"TestBot/src/models"
	"TestBot/src/periodic_task/distlock"

	"os"
	"time"
	"fmt"
)

type PeriodicFunc func(ctx *models.Context)

func distlockWrapper(operation PeriodicFunc, frequencyS int, lockName string, ctx *models.Context) {
	ticker := time.NewTicker(time.Duration(frequencyS) * time.Second)

	for range ticker.C {
		ownerID, err := os.Hostname()
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to gen hostname: %v", err))
			continue
		}

		locked, err := distlock.AcquireLock(ctx.Postgres(), lockName, ownerID)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to get lock: %v", err))
			continue
		}

		if locked {
			operation(ctx)

			err = distlock.ReleaseLock(ctx.Postgres(), lockName, ownerID)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to unclock: %v", err))
			}
		}
	}
}

func periodicTaskWrapper(operation PeriodicFunc, frequencyS int, ctx *models.Context) {
	ticker := time.NewTicker(time.Duration(frequencyS) * time.Second)

	for range ticker.C {
		operation(ctx)
	}
}

func Start(ctx *models.Context) {
	go distlockWrapper(SendReminders, ctx.Config.PeriodicTask.SendRemindersSettings.FrequencyS, ctx.Config.PeriodicTask.SendRemindersSettings.LockName, ctx)
	go distlockWrapper(DeleteExpiredReminder, ctx.Config.PeriodicTask.DeleteExpiredRemindersSettings.FrequencyS, ctx.Config.PeriodicTask.DeleteExpiredRemindersSettings.LockName, ctx)
	go periodicTaskWrapper(ClearExpiredState, ctx.Config.PeriodicTask.ClearExpiredStateSettings.FrequencyS, ctx)
	go periodicTaskWrapper(DumpToFile, ctx.Config.PeriodicTask.DumpToFileSettings.FrequencyS, ctx)
	go periodicTaskWrapper(LoadConfig, ctx.Config.PeriodicTask.LoadConfigSettings.FrequencyS, ctx)
}
