package periodictask

import(
	"TestBot/src/models"

	"time"
)

func ClearExpiredState(ctx *models.Context) {
	ctx.UserCache.ClearExpiredState(time.Duration(ctx.Config.PeriodicTask.ClearExpiredStateSettings.TtlS) * time.Second)
}
