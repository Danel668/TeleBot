package handlers

import(
	"gopkg.in/telebot.v4"

	"TestBot/src/caches"
	"TestBot/src/models"
	"TestBot/src/models/sources"
	"TestBot/src/data"
	"TestBot/src/utils"

	"regexp"
	"strconv"
	"strings"
	"fmt"
	"time"
)

const (
	Creator = "creator"
	Admin = "admin"
	RegexRec = `^rec\s\d+$`
	RegexBan = `^ban\s\d{1,19}\s(?:all|recommendation)\s.+$`
	RegexAssignRole = `^assign_role\s\d{1,19}\s(?:admin|user)$`
	RegexConfigGet = `^config\sget$`
	RegexConfigChange = `^config\s\S+\s.+$`
	RegexAddReminderTodayPattern = `^\d{2}:\d{2}\s.+$`
	RegexAddReminderCommonPattern = `^\d{2}\.\d{2}\.\d{4}\s\d{2}:\d{2}\s.+$`
	kToday = "сегодня"
	kYesterday = "вчера"
	kRegexGetRationsCommonPattern = `^\d{2}\.\d{2}\.\d{4}$`
)

func addRecommendationStateHandler(ctx *models.Context, userId int64, recommendation string) error {
	err := data.AddRecommendation(ctx.Postgres(), userId, recommendation)

	if err == nil {
		ctx.UserCache.SetState(userId, caches.Start)
	}

	return err
}

func startStateHandler(ctx *models.Context, userId int64, userInput string, c telebot.Context) error {
	if userInput == ctx.Config.Section.SectionAdmin.AdminModeSwitching {

		role, err := data.GetUserRoleByUserId(ctx.Postgres(), userId)

		if err != nil {
			return err
		}

		if role == Creator || role == Admin {
			ctx.UserCache.SetState(userId, caches.Admin)
			return c.Send(ctx.Config.Section.SectionAdmin.TextAfterSwitchingAdminMode)
		}
	}
	return nil
}

func adminStateHandler(ctx *models.Context, userInput string, c telebot.Context) error {
	userInput = strings.TrimSpace(userInput)

	regexRec := regexp.MustCompile(RegexRec)
	regexBan := regexp.MustCompile(RegexBan)
	regexAssingRole := regexp.MustCompile(RegexAssignRole)
	regexConfigGet := regexp.MustCompile(RegexConfigGet)
	regexConfigChange := regexp.MustCompile(RegexConfigChange)

	if regexRec.MatchString(userInput) {
		parts := strings.Split(userInput, " ")

		if len(parts) == 2 {
			limit, err := strconv.Atoi(parts[1])
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("adminStateHandler error: %v", err))
				return err
			}

			recommendations, err := data.GetRecommendations(ctx.Postgres(), limit)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("adminStateHandler error to get recommendations: %v", err))
				return err
			}

			var message string
			for _, rec := range recommendations {
				oneRec := fmt.Sprintf("UserId: %d\nRecommendation: %s\nSend_at: %s\n", rec.UserId, rec.Recommendation, rec.SendAt)

				message += oneRec
				message += "\n"
			}
			return c.Send(message)

		} else {
			ctx.Logger.Info("adminStateHandler parts != 2")
			return nil
		}
	} else if regexBan.MatchString(userInput) {
		parts := strings.Split(userInput, " ")

		if len(parts) == 4 {
			userId, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("adminStateHandler, error to parse userId: %v", err))
				return nil
			}
			bannedSection := parts[2]
			reason := parts[3]

			err = data.BanUser(ctx.Postgres(), userId, reason, bannedSection)

			if err != nil {
				return err
			}

			return c.Send(ctx.Config.Section.SectionAdmin.TextAfterBanUser)
		} else {
			ctx.Logger.Info("adminStateHandler, parts != 4")
			return nil
		}
	} else if regexAssingRole.MatchString(userInput) {
		parts := strings.Split(userInput, " ")

		if len(parts) == 3 {
			userId, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("adminStateHandler, error to parse userId: %v", err))
				return nil
			}
			role := parts[2]

			err = data.AssignRole(ctx.Postgres(), userId, role)
			if err != nil {
				return err
			}

			return c.Send(ctx.Config.Section.SectionAdmin.TextAfterAssignRole)

		} else {
			ctx.Logger.Info("adminStateHandler, parts != 3")
			return nil
		}
	} else if regexConfigGet.MatchString(userInput) {
		parts := strings.Split(userInput, " ")

		if len(parts) == 2 {
			configString, err := sources.GetConfigToString()
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to gen configString: %v", err))
				return nil
			}

			return c.Send(configString)
		} else {
			ctx.Logger.Info("adminStateHandler, parts != 2")
			return nil
		}
	} else if regexConfigChange.MatchString(userInput) {
		parts := strings.Split(userInput, " ")

		if len(parts) >= 3 {
			key := utils.ToGoName(parts[1])
			value := strings.Join(parts[2:], " ")
			
			err := sources.ChangeConfigValueByKey(key, value, ctx.Logger)

			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("adminStateHandler, failed to change config: %v", err))
				return nil
			}

			return c.Send(ctx.Config.Section.SectionAdmin.TextAfterChangingConfig)

		} else {
			ctx.Logger.Info("adminStateHandler, parts < 3")
			return nil
		}
	}

	return nil
}

func registrationTimeZoneStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {
	timezone, err := utils.ConvertTimeZone(userInput)
	if err != nil {
		return c.Send(ctx.Config.Section.SectionRegistration.TextAfterFailedRegistration)
	}

	err = data.Registration(ctx.Postgres(), userId, timezone)
	if err != nil {
		return err
	}

	ctx.UserCache.SetState(userId, caches.Start)
	return c.Send(ctx.Config.Section.SectionRegistration.TextAfterSuccessfulRegistration)
}

func countOfCaloriesStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {

	if len(userInput) > 100 {
		return c.Send(ctx.Config.Section.SectionCalories.TextIfRequestExceedsCondition)
	}

	gptAnswer := "Примерно 400 ккал."// идем с userInput в YA GPT (только чуть исправляем, добавляем по типу: посчитай калории userInput, и ограничения на токены)

	ctx.UserCache.SetState(userId, caches.Start)
	inlineKeys := [][]telebot.InlineButton{
		{ctx.Section.SectionCalories.InlineButtonChangeRequest},
	}

	return c.Send(gptAnswer, &telebot.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func getMyRationsStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {
	userInput = strings.TrimSpace(userInput)

	regexGetRationsCommonPattern := regexp.MustCompile(kRegexGetRationsCommonPattern)
	var userTime time.Time

	if userInput == kToday {
		userTime = time.Now()
	} else if userInput == kYesterday {
		userTime = time.Now().AddDate(0, 0, -1)
	} else if regexGetRationsCommonPattern.MatchString(userInput) {
		var err error
		userTime, err = utils.ParseDateInCommonFormat(userInput)
		if err != nil {
			return c.Send(ctx.Config.Section.SectionMyProgress.TextErrorFormatGetMyRations)
		}
	} else {
		return c.Send(ctx.Config.Section.SectionMyProgress.TextErrorFormatGetMyRations)
	}

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

	userTime = userTime.In(location)

	rations, err := data.GetUserRation(ctx.Postgres(), userId, userTime)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get user ration for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	ctx.UserCache.SetState(userId, caches.Start)
	if len(rations) == 0 {
		return c.Send(ctx.Config.Section.SectionMyProgress.TextEmptyRations)
	}
	return c.Send(strings.Join(rations, "\n"))
}

func addReminderStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {
	userInput = strings.TrimSpace(userInput)

	regexAddReminderTodayPattern := regexp.MustCompile(RegexAddReminderTodayPattern)
	regexAddReminderCommonPattern := regexp.MustCompile(RegexAddReminderCommonPattern)
	
	if regexAddReminderTodayPattern.MatchString(userInput) {
		parts := strings.Split(userInput, " ")
		if len(parts) >= 2 {
			timezone, err := data.GetTimezoneByPrimaryKey(ctx.Postgres(), userId)

			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to get timezone for user: %d, error: %v", userId, err))
				return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
			}

			timeString := parts[0]
			timeParsed, err := utils.ParseHoursMinsString(timeString, timezone)
			if err != nil {
				return c.Send(ctx.Config.Section.SectionReminder.TextInvalidPatternAddReminder)
			}

			localTime := timeParsed.Local()
			if !localTime.After(time.Now()) {
				return c.Send(ctx.Config.Section.SectionReminder.TextInvalidTime)
			}

			nameReminder := strings.Join(parts[1:], " ")
			err = data.InsertReminder(
				ctx.Postgres(), 
				userId, 
				localTime, 
				localTime.Add(time.Duration(ctx.Config.Section.SectionReminder.DifferenceBetweenSendingAndExpiringM) * time.Minute), 
				nameReminder,
			)

			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to insert reminder for user: %d, error: %v", userId, err))
				return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
			}

			ctx.UserCache.SetState(userId, caches.Start)
			return c.Send(ctx.Config.Section.SectionReminder.TextAfterAddingReminder)
		} else {
			return c.Send(ctx.Config.Section.SectionReminder.TextInvalidPatternAddReminder)
		}
	} else if regexAddReminderCommonPattern.MatchString(userInput) {
		parts := strings.Split(userInput, " ")
		if len(parts) >= 3 {
			timezone, err := data.GetTimezoneByPrimaryKey(ctx.Postgres(), userId)

			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to get timezone for user %d, error %v", userId, err))
				return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
			}

			timeString := strings.Join(parts[:2], " ")
			timeParsed, err := utils.ParseTimeInCommonFormat(timeString, timezone)
			if err != nil {
				return c.Send(ctx.Config.Section.SectionReminder.TextInvalidPatternAddReminder)
			}

			localTime := timeParsed.Local()
			if !localTime.After(time.Now()) {
				return c.Send(ctx.Config.Section.SectionReminder.TextInvalidTime)
			}

			nameReminder := strings.Join(parts[2:], " ")
			err = data.InsertReminder(
				ctx.Postgres(), 
				userId,
				localTime,
				localTime.Add(time.Duration(ctx.Config.Section.SectionReminder.DifferenceBetweenSendingAndExpiringM) * time.Minute),
				nameReminder,
			)

			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to insert reminder for user: %d, error: %v", userId, err))
				return c.Send(ctx.Config.Section.SectionReminder.TextErrorAfterCheckRegistration)
			}

			ctx.UserCache.SetState(userId, caches.Start)
			return c.Send(ctx.Config.Section.SectionReminder.TextAfterAddingReminder)
		} else {
			return c.Send(ctx.Config.Section.SectionReminder.TextInvalidPatternAddReminder)
		}
	}

	return c.Send(ctx.Config.Section.SectionReminder.TextInvalidPatternAddReminder)
}

func addRationStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {
	if len(userInput) > 100 {
		return c.Send(ctx.Config.Section.SectionCalories.TextIfRequestExceedsCondition)
	}

	goal, err := data.GetGoal(ctx.Postgres(), userId)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get goal for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	userTimeZone, err := data.GetTimezoneByPrimaryKey(ctx.Postgres(), userId)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get user timezone for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	now := time.Now()
	location, err := time.LoadLocation(userTimeZone)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get location for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	userTime := now.In(location)

	pastRations, err := data.GetUserRation(ctx.Postgres(), userId, userTime)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to get past ration for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}
	pastRations = append(pastRations, userInput)

	err = data.InsertRation(ctx.Postgres(), userId, userInput, userTime)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to add ration for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	// сформирвоать реквест в GPT и сходить
	answerGPT := fmt.Sprintf("Примерно 100 калорий %d %s", goal, strings.Join(pastRations, " "))

	ctx.UserCache.SetState(userId, caches.Start)
	return c.Send(answerGPT)
}

func setGoalStateHandler(ctx *models.Context, userInput string, userId int64, c telebot.Context) error {
	goal, err := strconv.Atoi(userInput)

	if err != nil {
		return c.Send(ctx.Config.Section.SectionMyProgress.TextErrorGoalFormat)
	}

	if goal <= 0 {
		return c.Send(ctx.Config.Section.SectionMyProgress.TextErrorGoalFormat)
	}

	err = data.InsertGoal(ctx.Postgres(), userId, goal)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("failed to set goal for user: %d, error: %v", userId, err))
		return c.Send(ctx.Config.Section.SectionMyProgress.TextError)
	}

	ctx.UserCache.SetState(userId, caches.Start)
	return c.Send(ctx.Config.Section.SectionMyProgress.TextAfterInsertingGoal)
}

func OnText(ctx *models.Context) func(c telebot.Context) error {
	return func(c telebot.Context) error {

		userId := c.Sender().ID
		state := ctx.UserCache.GetState(userId)

		switch state {
		case caches.AddRecommendation: {
			err := addRecommendationStateHandler(ctx, userId, c.Text())
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in addRecommendationStateHandler for user: %d, error: %v", userId, err))
				return c.Send(ctx.Config.Section.SectionRecommendation.TextErrorAfterGettingRecommendation)
			}
	
			return c.Send(ctx.Config.Section.SectionRecommendation.TextAfterGettingRecommendation)
		}
		case caches.Start: {
			err := startStateHandler(ctx, userId, c.Text(), c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in startStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.Admin: {
			err := adminStateHandler(ctx, c.Text(), c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in adminStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.RegistrationTimeZone: {
			err := registrationTimeZoneStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in registrationTimeZoneStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.CountOfCalories: {
			err := countOfCaloriesStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in countOfCaloriesStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.AddReminder: {
			err := addReminderStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in AddReminderStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.AddRation: {
			err := addRationStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in addRationStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.SetGoal: {
			err := setGoalStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in setGoalStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		case caches.GetMyRation: {
			err := getMyRationsStateHandler(ctx, c.Text(), userId, c)
			if err != nil {
				ctx.Logger.Error(fmt.Sprintf("error in getMyRationsStateHandler for user: %d, error: %v", userId, err))
			}
			return nil
		}
		default:
		}

		return nil
	}
}
