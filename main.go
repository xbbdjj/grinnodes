package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/xbbdjj/grinnodes/chart"
	"github.com/xbbdjj/grinnodes/client"
	"github.com/xbbdjj/grinnodes/config"
	"github.com/xbbdjj/grinnodes/ip"
	"github.com/xbbdjj/grinnodes/p2p"
	"github.com/xbbdjj/grinnodes/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	go client.Start()
	go client.Sync()
	go p2p.Start()
	go ip.Start()
	go chart.Start()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/googlemap", func(c *gin.Context) {

		gps, _ := storage.AllGPS()

		type country struct {
			Name    string
			Total   int
			Percent string
		}

		total, _ := storage.NodeTotal()
		publictotal, _ := storage.NodePublicCount()

		countrys := []country{}
		arr, _ := storage.AllCountry()
		for _, v := range arr {
			c := country{
				Name:    v.Name,
				Total:   v.Total,
				Percent: fmt.Sprintf("%.2f", float32(v.Total*100)/float32(total)),
			}
			countrys = append(countrys, c)
		}

		c.HTML(http.StatusOK, "indexgoogle.tmpl", gin.H{
			"country": countrys,
			"latlng":  gps,
			"total":   total,
			"public":  publictotal,
			"height":  client.Status.Height,
		})
	})

	r.GET("/", func(c *gin.Context) {
		type country struct {
			Name    string
			Total   int
			Percent string
		}

		total, _ := storage.NodeTotal()
		publictotal, _ := storage.NodePublicCount()

		countrys := []country{}
		arr, _ := storage.AllCountry()
		for _, v := range arr {
			c := country{
				Name:    v.Name,
				Total:   v.Total,
				Percent: fmt.Sprintf("%.2f", float32(v.Total*100)/float32(total)),
			}
			countrys = append(countrys, c)
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"country":   countrys,
			"total":     total,
			"public":    publictotal,
			"height":    client.Status.Height,
			"mapboxkey": config.Conf.MapBoxKey,
		})
	})

	r.GET("/geojson", func(c *gin.Context) {
		geo, _ := storage.GetGeoJSON()
		c.JSON(http.StatusOK, geo)
	})

	r.POST("/search", func(c *gin.Context) {
		ip := net.ParseIP(c.PostForm("ip"))
		if ip == nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"error": "please type right ip address!",
			})
			return
		}
		type node struct {
			IsPublic  bool
			Address   string
			RDNS      string
			UserAgent string
			Height    int
			LastSeen  string
			Country   string
			City      string
			NetWork   string
		}

		peers := []node{}

		res, _ := storage.Search(ip.String())
		if len(res) == 0 {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"error": "unable to found the node!",
			})
			return
		}
		for _, v := range res {
			n := node{
				Address:   fmt.Sprintf("%s:%d", v.IP, v.Port),
				UserAgent: v.UserAgent,
				Height:    v.Height,
				Country:   v.CountryName,
				City:      v.CityName,
				NetWork:   v.Org,
			}
			if v.RDNS != "" {
				n.RDNS = fmt.Sprintf("rDNS:%s", v.RDNS)
			}
			lastseen := 0
			if v.P2PLastConnected > lastseen {
				lastseen = v.P2PFirstConnected
			}
			if v.P2PLastSeen > lastseen {
				lastseen = v.P2PLastSeen
			}
			if v.APILastSeen > lastseen {
				lastseen = v.APILastSeen
			}
			n.LastSeen = time.Unix(int64(lastseen), 0).In(time.UTC).Format("2006-01-02 03:04:05")
			if int64(v.P2PLastConnected) > time.Now().Unix()-3600 {
				n.IsPublic = true
			}
			peers = append(peers, n)
		}

		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"peers": peers,
		})

	})

	r.GET("/nodes", func(c *gin.Context) {

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			page = 1
		}
		if page <= 0 {
			page = 1
		}

		type node struct {
			IsPublic  bool
			Address   string
			RDNS      string
			UserAgent string
			Height    int
			LastSeen  string
			Country   string
			City      string
			NetWork   string
		}

		peers := []node{}

		res, _ := storage.NodeList(page)
		for _, v := range res {
			n := node{
				Address:   fmt.Sprintf("%s:%d", v.IP, v.Port),
				UserAgent: v.UserAgent,
				Height:    v.Height,
				Country:   v.CountryName,
				City:      v.CityName,
				NetWork:   v.Org,
			}
			if v.RDNS != "" {
				n.RDNS = fmt.Sprintf("rDNS:%s", v.RDNS)
			}
			lastseen := 0
			if v.P2PLastConnected > lastseen {
				lastseen = v.P2PFirstConnected
			}
			if v.P2PLastSeen > lastseen {
				lastseen = v.P2PLastSeen
			}
			if v.APILastSeen > lastseen {
				lastseen = v.APILastSeen
			}
			n.LastSeen = time.Unix(int64(lastseen), 0).In(time.UTC).Format("2006-01-02 03:04:05")
			if int64(v.P2PLastConnected) > time.Now().Unix()-3600 {
				n.IsPublic = true
			}
			peers = append(peers, n)
		}

		count, err := storage.NodeTotal()
		pageTotal := count/20 + 1

		pageLink := []int{}
		if page-5 > 1 {
			pageLink = append(pageLink, 1)
		}
		if page-5 > 2 {
			pageLink = append(pageLink, 0)
		}
		for p := page - 5; p <= page+5; p++ {
			if p >= 1 && p <= pageTotal {
				pageLink = append(pageLink, p)
			}
		}
		if page+5 < pageTotal-1 {
			pageLink = append(pageLink, 0)
		}
		if page+5 < pageTotal {
			pageLink = append(pageLink, pageTotal)
		}

		c.HTML(http.StatusOK, "nodes.tmpl", gin.H{
			"peers":    peers,
			"page":     page,
			"pageLink": pageLink,
		})
	})

	r.GET("/history", func(c *gin.Context) {
		arr, _ := storage.GetChart()
		c.HTML(http.StatusOK, "highchart.tmpl", gin.H{
			"total": arr,
		})
	})

	r.Run(":80")
}
