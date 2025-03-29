package handlers

import(
	"gopkg.in/telebot.v4"

	"TestBot/src/caches"
	"TestBot/src/models"
)

func Recommendation(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.AddRecommendation)

		return c.Send(ctx.Config.Section.SectionRecommendation.TextAfterButtonRecommendation)
	}
}
