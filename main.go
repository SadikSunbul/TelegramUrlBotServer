package main

import (
	"github.com/SadikSunbul/TelegramUrlBotServer/config"
	Database "github.com/SadikSunbul/TelegramUrlBotServer/database"
	"github.com/SadikSunbul/TelegramUrlBotServer/fiber"
)

func main() {
	config.LoadConfig("config.yaml")
	Database.ConnectionDatabase()

	fiber.Root()
}
