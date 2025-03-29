package handlers

import(
	"gopkg.in/telebot.v4"

	"TestBot/src/caches"
	"TestBot/src/models"
	"TestBot/src/data"

	"fmt"
)

func Reminder(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {

		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.Start)

		isRegistrated, err := data.IsRegistrated(ctx.Postgres(), userId)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to check is user registrated: %v", err))
			return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
		}

		if !isRegistrated {
			handler := Registration(ctx)
			return handler(c)
		}

		keyboard := [][]telebot.ReplyButton{
			{
				ctx.Section.SectionReminder.ButtonAddReminder,
				ctx.Section.SectionReminder.ButtonReminderManual,
				ctx.Section.SectionReminder.ButtonGetMyReminders,
			},
		}

		return c.Send(ctx.Config.Section.SectionReminder.TextAfterButtonReminder, &telebot.ReplyMarkup{
			ReplyKeyboard:       keyboard,
			RemoveKeyboard: true,
			OneTimeKeyboard:     true,
		})
	}
}

func AddReminder(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.AddReminder)

		return c.Send(ctx.Config.Section.SectionReminder.TextAfterButtonAddReminder)
	}
}

func ReminderManual(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		return c.Send(ctx.Config.Section.SectionReminder.TextAfterButtonReminderManual)
	}
}

func GetMyReminder(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.Start)

		reminders, err := data.GetActualRemindersByUserId(ctx.Postgres(), userId)
		if err != nil {
			return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
		}

		if len(reminders) == 0 {
			return c.Send(ctx.Config.Section.SectionReminder.TextIfRemindersIsEmpty)
		}

		var message string
		for _, reminder := range reminders {
			message += fmt.Sprintf("%s: %s\n", reminder.SendAt, reminder.Reminder)
		}

		return c.Send(message)
	}
}
