package Database

import (
	"context"
	"github.com/SadikSunbul/TelegramUrlBotServer/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
	Client *mongo.Database
}

const (
	Url    string = "urls"
	UrlIfo string = "url_infos"
)

func ConnectionDatabase() *DataBase {
	cfg := config.GetConfig()
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(cfg.MongoDbConnect))
	if err != nil {
		panic(err)
	}
	return &DataBase{Client: client.Database(cfg.DbName)}
}
