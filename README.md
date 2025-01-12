# TelegramUrlBotServer

Telegram Url Bot Server

Bu proje, [TelegramUrlBot](https://github.com/SadikSunbul/TelegramUrlBot) için API sunucusu olarak çalışır. Bu API, kısaltılmış URL'lere erişim sağlar ve kullanıcıdan aşağıdaki verileri toplar:

### Kullanıcı Cihaz Bilgileri

```go
type UserDeviceInfo struct {
    IP           string             `json:"ip" bson:"ip"`
    UrlId        primitive.ObjectID `json:"urlId" bson:"urlId"`
    ClickTime    primitive.DateTime `json:"clickTime,omitempty" bson:"clickTime,omitempty"`
    Device       string             `json:"device" bson:"device"`
    Browser      string             `json:"browser" bson:"browser"`
    OS           string             `json:"os" bson:"os"`
    UserAgent    string             `json:"user_agent" bson:"user_agent"`
    Country      string             `json:"country" bson:"country"`
    CountryCode  string             `json:"country_code" bson:"country_code"`
    Region       string             `json:"region" bson:"region"`
    City         string             `json:"city" bson:"city"`
    ISP          string             `json:"isp" bson:"isp"`
    Organization string             `json:"organization" bson:"organization"`
    ASN          string             `json:"asn" bson:"asn"`
    Timezone     string             `json:"timezone" bson:"timezone"`
    Location     Location           `json:"location" bson:"location"`
}

type Location struct {
    Latitude  float64 `json:"latitude" bson:"latitude"`
    Longitude float64 `json:"longitude" bson:"longitude"`
}
