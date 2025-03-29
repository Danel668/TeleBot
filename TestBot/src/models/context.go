package models

import (
	"gopkg.in/telebot.v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"log"
	"context"
	"os"
	"time"
	"fmt"

	"TestBot/src/models/sources"
	"TestBot/src/utils/logger"
	"TestBot/src/caches"
)

type SectionRecommendation struct {
	ButtonRecommendation telebot.ReplyButton
}

type SectionRegistration struct {
	ButtonRegistration telebot.ReplyButton
}

type SectionCalories struct {
	ButtonCalories            telebot.ReplyButton
	InlineButtonChangeRequest telebot.InlineButton
}

type SectionReminder struct {
	ButtonReminder       telebot.ReplyButton
	ButtonAddReminder    telebot.ReplyButton
	ButtonReminderManual telebot.ReplyButton
	ButtonGetMyReminders telebot.ReplyButton
}

type SectionMyProgress struct {
    ButtonMyProgress                        telebot.ReplyButton
    ButtonMyProgressManual                  telebot.ReplyButton
    ButtonMyProgressSetGoal                 telebot.ReplyButton
    ButtonMyProgressAddRation               telebot.ReplyButton
    ButtonMyProgressGetMyRations            telebot.ReplyButton
	ButtonMyProgressGetMyRationsForLastTime telebot.ReplyButton
}

type Section struct {
	SectionRecommendation SectionRecommendation
	SectionRegistration   SectionRegistration
	SectionCalories       SectionCalories
	SectionReminder       SectionReminder
	SectionMyProgress     SectionMyProgress
}

type Context struct {
	Bot        *telebot.Bot
	Config     *sources.Config
	Section    Section
	DBPool     *pgxpool.Pool
	Logger     *zap.Logger
	FileLogger *os.File
	UserCache  *caches.StateUserCache
}

func(c *Context) Postgres() *pgxpool.Conn {
	conn, err := c.DBPool.Acquire(context.Background())
	if err != nil {
		c.Logger.Warn(fmt.Sprintf("failed to get connection: %v", err))
		return nil
	}
	return conn
}

func NewContext() *Context {

	connStrPostgres := os.Getenv("POSTGRES_CONNECTION_STRING")
	conf, err := pgxpool.ParseConfig(connStrPostgres)
	if err != nil {
		log.Fatalln("failed to parse config with db string:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		log.Fatalln("failed to get pool connection:", err)
	}

	fileLogger, err := os.OpenFile("./logs/production.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("failed to open logs file for writing logs")
	}
	logger := logger.NewLogger(fileLogger)

	userCache := caches.NewStateUserCache()

	config := sources.NewConfig(logger)
	if config == nil {
		log.Fatalln("failed to create the config")
	}

	bot, err := sources.TelebotInitializer()
	if err != nil {
		log.Fatalln("failed to initialize telebot:", err)
	}

	context := &Context{
        Bot: bot,
		Config: config,
		Section: Section{
			SectionRecommendation: SectionRecommendation{
				ButtonRecommendation: telebot.ReplyButton{
					Text: config.Section.SectionRecommendation.ButtonRecommendation,
				},
			},
			SectionRegistration: SectionRegistration{
				ButtonRegistration: telebot.ReplyButton{
					Text: config.Section.SectionRegistration.ButtonRegistration,
				},
			},
			SectionCalories: SectionCalories{
				ButtonCalories: telebot.ReplyButton{
					Text: config.Section.SectionCalories.ButtonCalories,
				},
				InlineButtonChangeRequest: telebot.InlineButton{
					Unique: config.Section.SectionCalories.InlineButtonChangeRequestUnique,
					Text: config.Section.SectionCalories.InlineButtonChangeRequest,
				},
			},
			SectionReminder: SectionReminder{
				ButtonReminder: telebot.ReplyButton{
					Text: config.Section.SectionReminder.ButtonReminder,
				},
				ButtonAddReminder: telebot.ReplyButton{
					Text: config.Section.SectionReminder.ButtonAddReminder,
				},
				ButtonReminderManual: telebot.ReplyButton{
					Text: config.Section.SectionReminder.ButtonReminderManual,
				},
				ButtonGetMyReminders:  telebot.ReplyButton{
					Text: config.Section.SectionReminder.ButtonGetMyReminders,
				},
			},
			SectionMyProgress: SectionMyProgress{
				ButtonMyProgress: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgress,
				},
				ButtonMyProgressManual: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgressManual,
				},
				ButtonMyProgressSetGoal: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgressSetGoal,
				},
				ButtonMyProgressAddRation: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgressAddRation,
				},
				ButtonMyProgressGetMyRations: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgressGetMyRations,
				},
				ButtonMyProgressGetMyRationsForLastTime: telebot.ReplyButton{
					Text: config.Section.SectionMyProgress.ButtonMyProgressGetMyRationsForLastTime,
				},
			},
		},
		DBPool: pool,
		Logger: logger,
		FileLogger: fileLogger,
		UserCache: userCache,
    }

	return context
}
