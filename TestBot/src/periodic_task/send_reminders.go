package periodictask

import(
	"TestBot/src/models"
	"TestBot/src/data"

	"gopkg.in/telebot.v4"

	"fmt"
)

func SendReminders(ctx *models.Context) {
	reminders, err := data.GetActualReminders(ctx.Postgres())
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get reminbers: %v", err))
		return
	}

	var sentReminders []data.Reminder
	for _, reminder := range reminders {
		userId := reminder.UserId
		message := fmt.Sprintf(ctx.Config.PeriodicTask.SendRemindersSettings.AnswerPlaceholder, reminder.Reminder)

		recipient := telebot.Recipient(&telebot.User{ID: userId})
		_, err := ctx.Bot.Send(recipient, message)

		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to send reminder to user: %d, err: %v", userId, err))
		} else {
			sentReminders = append(sentReminders, reminder)
		}
	}

	if len(sentReminders) > 0 {
		err = data.DeleteSentReminders(ctx.Postgres(), sentReminders)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to delete sent reminders: %v", err))
		}
	}
}
