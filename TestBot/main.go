package main

import (
	"TestBot/src/models"
	"TestBot/src/utils/initializers"
	"TestBot/src/middleware"
	"TestBot/src/periodic_task"
)

func main() {
	context := models.NewContext()

	defer func() {
		context.UserCache.DumpToFile(context.Logger)
		context.DBPool.Close()
		context.Logger.Sync()
		context.FileLogger.Close()
	}()

	if context != nil {
		periodictask.Start(context)

		context.Bot.Use(middleware.TelebotMiddleware(context))
		initializers.TelebotHandlersInitializer(context)
		context.Bot.Start()
	}
}
