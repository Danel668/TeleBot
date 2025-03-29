package periodictask

import(
	"TestBot/src/models"
	"TestBot/src/models/sources"
)

func LoadConfig(ctx *models.Context) {
	config := sources.LoadConfig(ctx.Logger)
	if config != nil {
		ctx.Config = config
	}
}
