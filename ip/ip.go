package ip

import (
	"encoding/json"
	"fmt"
	"github.com/xbbdjj/grinnodes/config"
	"github.com/xbbdjj/grinnodes/log"
	"github.com/xbbdjj/grinnodes/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

//data from https://ipapi.co/
type ipapi struct {
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	City          string  `json:"city"`
	Asn           string  `json:"asn"`
	Org           string  `json:"org"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `json:"country"`
	CountryName   string  `json:"country_name"`
}

//data from https://ipdata.co/
type ipdataAsn struct {
	Asn    string `json:"asn"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Route  string `json:"route"`
	Type   string `json:"type"`
}

type ipdata struct {
	ContinentCode string    `json:"continent_code"`
	CountryCode   string    `json:"country_code"`
	CountryName   string    `json:"country_name"`
	Latitude      float32   `json:"latitude"`
	Longitude     float32   `json:"longitude"`
	City          string    `json:"city"`
	Asn           ipdataAsn `json:"asn"`
}

var ipapitime int

/*
state
0 init
ipapi.co -1 1
ipdata.co -2 2
*/

func Start() {
	time1 := time.Now()
	time2 := time.Now()
	for {
		time.Sleep(time.Minute)
		if time.Now().After(time1) {
			peers, err := storage.AllIP(0)
			if err != nil {
				log.Logger.Error("ipapi sql get error", zap.Error(err))
				continue
			}
			for _, ip := range peers {
				url := fmt.Sprintf("https://ipapi.co/%s/json", ip)
				resp, err := http.Get(url)

				if err != nil {
					continue
				}

				if resp.StatusCode == 429 {
					time1 = time.Now().Add(time.Hour * 24)
					break
				}

				var res ipapi
				json.NewDecoder(resp.Body).Decode(&res)
				defer resp.Body.Close()

				if res.Latitude == 0 && res.Longitude == 0 {
					storage.UpdateIPState(ip, -1)
					continue
				}

				info := storage.IPInfo{
					ContinentCode: res.ContinentCode,
					CountryCode:   res.CountryCode,
					CountryName:   res.CountryName,
					City:          res.City,
					Latitude:      res.Latitude,
					Longitude:     res.Longitude,
					Asn:           res.Asn,
					Org:           res.Org,
				}

				storage.UpdateIP(ip, info, 1)
			}
		}

		if time.Now().After(time2) {
			peers, err := storage.AllIP(-1)
			if err != nil {
				continue
			}
			for _, ip := range peers {

				key := config.NewConfig().IPDataKey
				url := fmt.Sprintf("https://api.ipdata.co/%s?api-key=%s", ip, key)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if resp.StatusCode == 403 {
					time2 = time.Now().Add(time.Hour * 24)
					break
				}

				var res ipdata
				json.NewDecoder(resp.Body).Decode(&res)
				defer resp.Body.Close()

				if res.Latitude == 0 && res.Longitude == 0 {
					storage.UpdateIPState(ip, -2)
					continue
				}

				info := storage.IPInfo{
					ContinentCode: res.ContinentCode,
					CountryCode:   res.CountryCode,
					CountryName:   res.CountryName,
					City:          res.City,
					Latitude:      res.Latitude,
					Longitude:     res.Longitude,
					Asn:           res.Asn.Asn,
					Org:           res.Asn.Name,
				}
				storage.UpdateIP(ip, info, 2)
			}
		}
	}
}
