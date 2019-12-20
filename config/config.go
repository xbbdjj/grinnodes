package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	GoogleMapKey       string
	MapBoxKey          string
	IPDataKey          string
	PeerActiveDuration int64
	ClientHttpAddr     string
}

var Conf config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./data/")
	viper.SetConfigType("toml")

	viper.SetDefault("googlemap_key", "")
	viper.SetDefault("mapbox_key", "")
	viper.SetDefault("ipdata_key", "")
	viper.SetDefault("peer_active_duration", 3600)
	viper.SetDefault("client_http_addr", "45.77.206.185:3413")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	Conf.GoogleMapKey = viper.GetString("googlemap_key")
	Conf.MapBoxKey = viper.GetString("mapbox_key")
	Conf.IPDataKey = viper.GetString("ipdata_key")
	Conf.PeerActiveDuration = viper.GetInt64("peer_active_duration")
	Conf.ClientHttpAddr = viper.GetString("client_http_addr")
}
