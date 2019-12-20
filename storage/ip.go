package storage

import (
	"fmt"
	"net"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//AllIP ... all ip to query the location
func AllIP(state int) ([]string, error) {
	ips := []string{}
	rows, err := db.Query("SELECT ip FROM peer WHERE ip_state = ? ORDER BY p2p_last_connected DESC", state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ips, nil
}

//UpdateIPState ... ip location request state
func UpdateIPState(ip string, state int) error {
	if _, err := db.Exec("UPDATE peer SET ip_state = ? WHERE ip = ? ", state, ip); err != nil {
		return err
	}
	return nil
}

//IPInfo ... ip info to store
type IPInfo struct {
	Latitude      float32
	Longitude     float32
	City          string
	Asn           string
	Org           string
	ContinentCode string
	CountryCode   string
	CountryName   string
}

//UpdateIP ... update ip info
func UpdateIP(ip string, info IPInfo, state int) error {
	rdns := ""
	a, err := net.LookupAddr(ip)
	if err == nil && len(a) >= 1 {
		rdns = strings.TrimRight(a[0], ".")
	}
	if _, err := db.Exec(
		"UPDATE peer SET ip_continent_code = ?, ip_country_code = ?, "+
			"ip_country_name = ?, ip_city = ?, ip_latitude = ?, "+
			"ip_longitude = ?, ip_asn = ?, ip_org = ?, ip_state = ?, ip_rdns = ? "+
			"WHERE ip = ? ",
		info.ContinentCode,
		info.CountryCode,
		info.CountryName,
		info.City,
		fmt.Sprintf("%.4f", info.Latitude),
		fmt.Sprintf("%.4f", info.Longitude),
		info.Asn,
		info.Org,
		state,
		rdns,
		ip,
	); err != nil {
		return err
	}
	return nil
}
