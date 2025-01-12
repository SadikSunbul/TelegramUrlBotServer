package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SadikSunbul/TelegramUrlBotServer/Models"
	Database "github.com/SadikSunbul/TelegramUrlBotServer/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IpApiResponse struct {
	Country string `json:"country"`
	Status  string `json:"status"`
}

func (db *UrlHandlers) Forward(ctx *fiber.Ctx) error {
	url := ctx.Params("url")
	data, err := db.DB.GetBy(Database.Url, bson.D{{"shortUrl", url}})
	if err != nil {
		ctx.SendString("Eror Database: " + err.Error())
		return err
	}
	get_url := Models.Url{}
	err = data.Decode(&get_url)
	if err != nil {
		ctx.SendString("Eror decode: " + err.Error())
		return err
	}

	if get_url.EndDate != 0 {
		currentTime := time.Now().Unix()
		if get_url.EndDate.Time().Unix() < currentTime {
			ctx.SendString("Bu urlnin kullanım zamanı dolmuştur.")
			return fiber.NewError(fiber.StatusGone, "URL süresi dolmuş")
		}
	}
	err = ctx.Redirect(get_url.OriginalUrl, fiber.StatusTemporaryRedirect)
	if err != nil {
		ctx.SendString("Bir hata oluştu, yönlendirme yapılamadı..")
		return err
	}

	var urlinfo Models.UrlInfo

	urlinfo.UrlId = get_url.Id
	urlinfo.ClickTime = primitive.NewDateTimeFromTime(time.Now())

	clientIP := ctx.IP()
	country := getCountryFromIP(clientIP)
	urlinfo.Country = country

	_, err = db.DB.Add(Database.UrlIfo, urlinfo)
	if err != nil {
		return err
	}
	return nil
}

func getCountryFromIP(ip string) string {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country", ip))
	if err != nil {
		return "Unknown"
	}
	fmt.Print("ip:", ip)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Unknown"
	}

	var result IpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Unknown"
	}

	if result.Status != "success" {
		return "Unknown"
	}
	return result.Country
}

func (db *UrlHandlers) IP(ctx *fiber.Ctx) error {
	// Gerçek IP adresini al
	clientIP := ctx.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = ctx.Get("X-Real-IP")
		if clientIP == "" {
			clientIP = ctx.IP()
		}
	}

	// Kullanıcının cihaz bilgilerini al
	userAgent := ctx.Get("User-Agent")
	device := getDeviceFromUserAgent(userAgent)

	// IP ve cihaz bilgilerini döndür
	response := fmt.Sprintf("IP Adresi: %s, Cihaz: %s", clientIP, device)
	ctx.SendString(response)
	return nil
}

// User-Agent'a göre cihazı belirle
func getDeviceFromUserAgent(userAgent string) string {
	if contains(userAgent, "Windows") {
		return "Windows"
	} else if contains(userAgent, "Macintosh") {
		return "Mac"
	} else if contains(userAgent, "iPhone") || contains(userAgent, "iPad") {
		return "iOS"
	} else if contains(userAgent, "Android") {
		return "Android"
	}
	return "Bilinmiyor"
}

// String içinde bir alt dizeyi kontrol et
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

/*
kullanıcın hangı ipv4 den ıstek attıgını bulur snra hangi cihazdan mac mı vındows mu ios mu yoksa androitmi oldugunu bulur sonra bulabildiği tüm verileri bulru ve kullanıcıya döner
*/
