package periodictask

import(
	"TestBot/src/models"
)

func DumpToFile(ctx *models.Context) {
	ctx.UserCache.DumpToFile(ctx.Logger)
}
