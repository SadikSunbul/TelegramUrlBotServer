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

	DatabaseRecord(ctx, db, get_url.Id)

	return nil
}

func DatabaseRecord(ctx *fiber.Ctx, db *UrlHandlers, urlId primitive.ObjectID) error {
	userInfo := Models.UserDeviceInfo{}

	// IP adresini al
	userInfo.IP = ctx.Get("X-Real-IP")
	if userInfo.IP == "" {
		userInfo.IP = ctx.Get("X-Forwarded-For")
		if userInfo.IP == "" {
			userInfo.IP = ctx.IP()
			if userInfo.IP == "" || userInfo.IP == "127.0.0.1" || userInfo.IP == "::1" {
				// Eğer localhost ise, gerçek IP'yi bulmak için harici servis kullan
				resp, err := http.Get("https://api.ipify.org?format=json")
				if err == nil {
					defer resp.Body.Close()
					var result struct {
						IP string `json:"ip"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
						userInfo.IP = result.IP
					}
				}
			}
		}
	}

	// User-Agent bilgilerini al
	userAgent := ctx.Get("User-Agent")
	userInfo.UserAgent = userAgent

	// Cihaz ve işletim sistemi tespiti
	if strings.Contains(userAgent, "Windows") {
		userInfo.OS = "Windows"
		if strings.Contains(userAgent, "Windows NT 10.0") {
			userInfo.OS += " 10"
		} else if strings.Contains(userAgent, "Windows NT 6.3") {
			userInfo.OS += " 8.1"
		} else if strings.Contains(userAgent, "Windows NT 6.2") {
			userInfo.OS += " 8"
		} else if strings.Contains(userAgent, "Windows NT 6.1") {
			userInfo.OS += " 7"
		}
	} else if strings.Contains(userAgent, "Macintosh") {
		userInfo.OS = "macOS"
	} else if strings.Contains(userAgent, "Linux") {
		if strings.Contains(userAgent, "Android") {
			userInfo.OS = "Android"
		} else {
			userInfo.OS = "Linux"
		}
	} else if strings.Contains(userAgent, "iPhone") {
		userInfo.OS = "iOS"
		userInfo.Device = "iPhone"
	} else if strings.Contains(userAgent, "iPad") {
		userInfo.OS = "iOS"
		userInfo.Device = "iPad"
	} else if strings.Contains(userAgent, "Android") {
		userInfo.OS = "Android"
	}

	// Tarayıcı tespiti
	switch {
	case strings.Contains(userAgent, "Chrome") && !strings.Contains(userAgent, "Edg") && !strings.Contains(userAgent, "OPR"):
		userInfo.Browser = "Chrome"
	case strings.Contains(userAgent, "Firefox"):
		userInfo.Browser = "Firefox"
	case strings.Contains(userAgent, "Safari") && !strings.Contains(userAgent, "Chrome"):
		userInfo.Browser = "Safari"
	case strings.Contains(userAgent, "Edg"):
		userInfo.Browser = "Edge"
	case strings.Contains(userAgent, "OPR"):
		userInfo.Browser = "Opera"
	default:
		userInfo.Browser = "Diğer"
	}

	// Cihaz tipini belirle
	if userInfo.Device == "" {
		if strings.Contains(userAgent, "Mobile") {
			userInfo.Device = "Mobil Cihaz"
		} else {
			userInfo.Device = "Masaüstü/Dizüstü"
		}
	}

	// IP detaylarını al (ipapi.co servisini kullanalım)
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,region,regionName,city,district,zip,lat,lon,timezone,isp,org,as,query", userInfo.IP))
	if err == nil {
		defer resp.Body.Close()
		var ipDetails Models.IpApiDetailResponse
		if err := json.NewDecoder(resp.Body).Decode(&ipDetails); err == nil && ipDetails.Status == "success" {

			userInfo.Country = ipDetails.Country
			userInfo.CountryCode = ipDetails.CountryCode
			userInfo.Region = ipDetails.RegionName
			userInfo.City = ipDetails.City
			userInfo.ISP = ipDetails.ISP
			userInfo.Organization = ipDetails.Org
			userInfo.ASN = ipDetails.AS
			userInfo.Timezone = ipDetails.Timezone
			userInfo.Location = Models.Location{
				Latitude:  ipDetails.Lat,
				Longitude: ipDetails.Lon,
			}
		} else {
			fmt.Printf("IP API Hata: %v\n", err)
		}
	} else {
		fmt.Printf("HTTP İstek Hatası: %v\n", err)
	}

	if userInfo.City == "" {
		userInfo.City = "Bilinmiyor"
	}

	userInfo.UrlId = urlId
	userInfo.ClickTime = primitive.NewDateTimeFromTime(time.Now())

	db.DB.Add(Database.UrlIfo, userInfo)

	return nil

}
