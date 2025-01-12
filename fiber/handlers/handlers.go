package handlers

import (
	Database "github.com/SadikSunbul/TelegramUrlBotServer/database"
)

type UrlHandlers struct {
	DB *Database.DataBase
}

func CreateUrlHandler() *UrlHandlers {
	return &UrlHandlers{
		DB: Database.ConnectionDatabase(),
	}
}
