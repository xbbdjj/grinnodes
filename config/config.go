package config

import (
	"github.com/spf13/viper"
)

type config struct {
	GoogleMapKey       string
	MapBoxKey          string
	IPDataKey          string
	PeerActiveDuration int64
	PeerClearDuration  int64
}

func NewConfig() config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./data/")
	viper.SetConfigType("toml")

	viper.SetDefault("googlemap_key", "")
	viper.SetDefault("mapbox_key", "")
	viper.SetDefault("ipdata_key", "")
	viper.SetDefault("peer_active_duration", 3600)
	viper.SetDefault("peer_clear_duration", 86400)

	viper.ReadInConfig()

	var conf config
	conf.GoogleMapKey = viper.GetString("googlemap_key")
	conf.MapBoxKey = viper.GetString("mapbox_key")
	conf.IPDataKey = viper.GetString("ipdata_key")
	conf.PeerActiveDuration = viper.GetInt64("peer_active_duration")
	conf.PeerClearDuration = viper.GetInt64("peer_clear_duration")

	return conf
}
