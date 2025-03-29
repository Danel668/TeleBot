package handlers

import(
	"gopkg.in/telebot.v4"

	"TestBot/src/caches"
	"TestBot/src/models"
)

func Start(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.Start)

		keyboard := [][]telebot.ReplyButton{
			{
				ctx.Section.SectionRecommendation.ButtonRecommendation,
				ctx.Section.SectionRegistration.ButtonRegistration,
				ctx.Section.SectionCalories.ButtonCalories,
				ctx.Section.SectionReminder.ButtonReminder,
				ctx.Section.SectionMyProgress.ButtonMyProgress,
			},
		}

		return c.Send("Бот запущен", &telebot.ReplyMarkup{
			ReplyKeyboard:       keyboard,
			RemoveKeyboard: true,
			OneTimeKeyboard:     true,
		})
	}
}
