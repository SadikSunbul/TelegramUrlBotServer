package fiber

import (
	"github.com/SadikSunbul/TelegramUrlBotServer/config"
	"github.com/SadikSunbul/TelegramUrlBotServer/fiber/handlers"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/log"
)

func Root() {
	app := fiber.New()
	cfg := config.GetConfig()
	//  Rootlar yazılıcak burada

	urlHandlers := handlers.CreateUrlHandler()

	// Önce sabit route'lar
	app.Get("/", func(ctx *fiber.Ctx) error {
		ctx.SendString("Telegram Url Bot Server active")
		return nil
	})
	// En son wildcard route
	app.Get("/:url", urlHandlers.Forward)

	log.Fatal(app.Listen(cfg.LolalHostPort))
}
