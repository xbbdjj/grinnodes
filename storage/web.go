package storage

import (
	"fmt"
	"time"

	"github.com/xbbdjj/grinnodes/config"

	_ "github.com/mattn/go-sqlite3"
)

type GpsInfo struct {
	Latitude  float32
	Longitude float32
}

func AllGPS() ([]GpsInfo, error) {
	t := time.Now().Unix() - config.NewConfig().PeerActiveDuration
	ips := []GpsInfo{}
	rows, err := db.Query(
		"SELECT ip_latitude, ip_longitude FROM peer WHERE ip_state > 0 AND "+
			"(p2p_last_connected > ? OR p2p_last_seen > ? OR api_last_seen > ?)",
		t,
		t,
		t,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var lat float32
		var lng float32
		err := rows.Scan(&lat, &lng)
		if err != nil {
			return nil, err
		}
		i := GpsInfo{
			Latitude:  lat,
			Longitude: lng,
		}
		ips = append(ips, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ips, nil
}

type CountryInfo struct {
	Name  string
	Total int
}

func AllCountry() ([]CountryInfo, error) {
	arr := []CountryInfo{}
	t := time.Now().Unix() - config.NewConfig().PeerActiveDuration
	rows, err := db.Query(
		"select ip_country_name, count(*) as t from peer where ip_country_name != '' AND "+
			" (p2p_last_connected > ? OR p2p_last_seen > ? OR api_last_seen > ?) "+
			" GROUP BY ip_country_name order by t desc",
		t,
		t,
		t,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var name string
		var total int
		err := rows.Scan(&name, &total)
		if err != nil {
			return nil, err
		}
		i := CountryInfo{
			Name:  name,
			Total: total,
		}
		arr = append(arr, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func NodeTotal() (int, error) {
	t := time.Now().Unix() - config.NewConfig().PeerActiveDuration
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM peer WHERE p2p_last_connected > ? OR p2p_last_seen > ? OR api_last_seen > ? ",
		t,
		t,
		t,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NodePublicCount() (int, error) {
	t := time.Now().Unix() - config.NewConfig().PeerActiveDuration
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM peer WHERE p2p_last_connected > ? ", t).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type PublicNode struct {
	IsActive  bool
	Address   string
	UserAgent string
	Height    int
	LastSeen  string
	Country   string
	City      string
	NetWork   string
}

func PublicNodeList(ip string) ([]PublicNode, error) {
	arr := []PublicNode{}
	sql := "SELECT ip, port, node_user_agent, node_height, " +
		"p2p_first_connected, p2p_last_connected, p2p_last_seen, api_last_seen, " +
		"ip_country_name, ip_city, ip_org, ip_rdns " +
		"FROM peer WHERE p2p_last_connected > 0 "
	if len(ip) > 0 {
		sql += " AND ip LIKE '%" + ip + "%' "
	}
	sql += "ORDER BY p2p_last_connected DESC, node_height DESC "

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

		var ip string
		var port int
		var userAgent string
		var height int
		var p2pFirstConnected int
		var p2pLastConnected int
		var p2pLastSeen int
		var apiLastSeen int
		var countryName string
		var cityName string
		var org string
		var rdns string
		err := rows.Scan(&ip, &port, &userAgent, &height, &p2pFirstConnected, &p2pLastConnected, &p2pLastSeen, &apiLastSeen, &countryName, &cityName, &org, &rdns)
		if err != nil {
			return nil, err
		}
		n := PublicNode{
			Address:   fmt.Sprintf("%s:%d", ip, port),
			UserAgent: userAgent,
			Height:    height,
			Country:   countryName,
			City:      cityName,
			NetWork:   org,
		}
		lastseen := p2pLastConnected
		if p2pLastSeen > p2pLastConnected {
			lastseen = p2pLastConnected
		}
		n.LastSeen = time.Unix(int64(lastseen), 0).In(time.UTC).Format("2006-01-02 03:04")

		if int64(p2pLastConnected) > time.Now().Unix()-config.NewConfig().PeerActiveDuration {
			n.IsActive = true
		}
		arr = append(arr, n)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

type geosjon struct {
	T        string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	T        string   `json:"type"`
	Geometry geometry `json:"geometry"`
}

type geometry struct {
	T           string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

func GetGeoJSON() (geosjon, error) {
	t := time.Now().Unix() - config.NewConfig().PeerActiveDuration
	geo := geosjon{T: "FeatureCollection"}
	rows, err := db.Query(
		"SELECT ip_latitude, ip_longitude FROM peer WHERE ip_state > 0 AND "+
			"(p2p_last_connected > ? OR p2p_last_seen > ? OR api_last_seen > ?)",
		t,
		t,
		t,
	)
	if err != nil {
		return geo, err
	}
	defer rows.Close()
	for rows.Next() {
		var lat float32
		var lng float32
		err := rows.Scan(&lat, &lng)
		if err != nil {
			return geo, err
		}
		g := geometry{
			T:           "Point",
			Coordinates: []float32{lng, lat},
		}
		f := feature{
			T:        "Feature",
			Geometry: g,
		}
		geo.Features = append(geo.Features, f)
	}
	err = rows.Err()
	if err != nil {
		return geo, err
	}
	return geo, nil
}
