package storage

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("sqlite3", "./data/node.db?cache=shared")
	db.SetMaxOpenConns(1)
}

func Test() {
	if _, err := db.Exec("INSERT INTO peer (ip) VALUES (?) ", "45.77.206.185"); err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
	}
}

func IsPublicIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 10 {
			return false
		} else if ip4[0] == 192 && ip4[1] == 168 {
			return false
		} else if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return false
		} else if ip4.String() == "127.0.0.1" {
			return false
		} else {
			return true
		}
	}
	return false
}
