package middleware

import (
	"TestBot/src/models"
	"TestBot/src/data"

	"gopkg.in/telebot.v4"

	"fmt"
)

func TelebotMiddleware(ctx *models.Context) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {

			userId := c.Sender().ID

			isBanned, err := data.IsUserBannedAll(ctx.Postgres(), userId)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to check banned_user: %v", err))
				return nil
			}

			if isBanned {
				return nil
			}

			return next(c)
		}
	}
}
