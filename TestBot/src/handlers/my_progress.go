package handlers

import(
	"gopkg.in/telebot.v4"

	"TestBot/src/caches"
	"TestBot/src/models"
	"TestBot/src/data"

	"time"
	"fmt"
)

func MyProgress(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.Start)

		isRegistrated, err := data.IsRegistrated(ctx.Postgres(), userId)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to check is user registrated for user: %d, error: %v", userId, err))
			return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
		}

		if !isRegistrated {
			handler := Registration(ctx)
			return handler(c)
		}

		keyboard := [][]telebot.ReplyButton{
			{
				ctx.Section.SectionMyProgress.ButtonMyProgressManual,
				ctx.Section.SectionMyProgress.ButtonMyProgressSetGoal,
				ctx.Section.SectionMyProgress.ButtonMyProgressAddRation,
				ctx.Section.SectionMyProgress.ButtonMyProgressGetMyRations,
			},
		}

		return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterButtonMyProgress, &telebot.ReplyMarkup{
			ReplyKeyboard:       keyboard,
			RemoveKeyboard: true,
			OneTimeKeyboard:     true,
		})
	}
}

func MyProgressManual(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.Start)

		return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterButtonMyProgressManual)
	}
}

func MyProgressAddRation(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.AddRation)

		return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterButtonAddRation)
	}
}

func MyProgressSetGoal(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.SetGoal)

		return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterButtonMyProgressSetGoal)
	}
}

func MyProgressGetMyRations(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID
		ctx.UserCache.SetState(userId, caches.GetMyRation)

		keyboard := [][]telebot.ReplyButton{
			{
				ctx.Section.SectionMyProgress.ButtonMyProgressGetMyRationsForLastTime,
			},
		}

		return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterButtonMyProgressGetMyRations, &telebot.ReplyMarkup{
			ReplyKeyboard:   keyboard,
			RemoveKeyboard:  true,
			OneTimeKeyboard: true,
		})
	}
}

func GetRationsForTheLastTime(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userId := c.Sender().ID


		now := time.Now()
		userTimeZone, err := data.GetTimezoneByPrimaryKey(ctx.Postgres(), userId)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to get user timezone for user: %d, error: %v", userId, err))
			return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
		}

		location, err := time.LoadLocation(userTimeZone)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to get location for user: %d, error: %v", userId, err))
			return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
		}

		userTime := now.In(location)
		startPeriod := userTime.AddDate(0, 0, -ctx.Config.Section.SectionMyProgress.LastTimeGetRations)
		
		rations, err := data.GetRationsForTheLastTime(ctx.Postgres(), userId, startPeriod)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to get rations for user: %d, error: %v", userId, err))
			return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
		}

		if len(rations) == 0 {
			return c.Send(ctx.Config.Section.SectionMyProgress.TextEmptyRations)
		}

		var answer string
		var lastDate string
		for _, ration := range rations {
			date := ration.CreatedAt.Format("02.01.2006")
			if date != lastDate {
				lastDate = date
				answer += fmt.Sprintf("\n%s:\n", lastDate)
			}
			answer += fmt.Sprintf("%s\n", ration.Ration)
		}
		
		ctx.UserCache.SetState(userId, caches.Start)
		return c.Send(answer)
	}
}
