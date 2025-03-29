package periodictask

import(
	"TestBot/src/models"
	"TestBot/src/data"

	"fmt"
)

func DeleteExpiredReminder(ctx *models.Context) {

	err := data.DeleteExpiredReminders(ctx.Postgres())
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to delete expired reminders: %v", err))
	}
}
