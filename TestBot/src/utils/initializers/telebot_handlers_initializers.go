package initializers

import (
	"TestBot/src/handlers"
	"TestBot/src/models"

	"gopkg.in/telebot.v4"
)

func TelebotHandlersInitializer(c *models.Context) {

	c.Bot.Handle("/start", handlers.Start(c))
	c.Bot.Handle(&c.Section.SectionRecommendation.ButtonRecommendation, handlers.Recommendation(c))
	c.Bot.Handle(telebot.OnText, handlers.OnText(c))
	c.Bot.Handle(&c.Section.SectionRegistration.ButtonRegistration, handlers.Registration(c))
	c.Bot.Handle(&c.Section.SectionCalories.ButtonCalories, handlers.CalculateCalories(c))
	c.Bot.Handle(&c.Section.SectionCalories.InlineButtonChangeRequest, handlers.CalculateCalories(c))
	c.Bot.Handle(&c.Section.SectionReminder.ButtonReminder, handlers.Reminder(c))
	c.Bot.Handle(&c.Section.SectionReminder.ButtonAddReminder, handlers.AddReminder(c))
	c.Bot.Handle(&c.Section.SectionReminder.ButtonReminderManual, handlers.ReminderManual(c))
	c.Bot.Handle(&c.Section.SectionReminder.ButtonGetMyReminders, handlers.GetMyReminder(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgress, handlers.MyProgress(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgressManual, handlers.MyProgressManual(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgressAddRation, handlers.MyProgressAddRation(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgressSetGoal, handlers.MyProgressSetGoal(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgressGetMyRations, handlers.MyProgressGetMyRations(c))
	c.Bot.Handle(&c.Section.SectionMyProgress.ButtonMyProgressGetMyRationsForLastTime, handlers.GetRationsForTheLastTime(c))
}
