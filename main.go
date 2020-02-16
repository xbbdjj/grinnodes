package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/xbbdjj/grinnodes/chart"
	"github.com/xbbdjj/grinnodes/config"

	"github.com/xbbdjj/grinnodes/ip"
	"github.com/xbbdjj/grinnodes/p2p"
	"github.com/xbbdjj/grinnodes/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	storage.ClearOldPeer()
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
			"mapboxkey": config.NewConfig().MapBoxKey,
		})
	})

	r.GET("/geojson", func(c *gin.Context) {
		geo, _ := storage.GetGeoJSON()
		c.JSON(http.StatusOK, geo)
	})

	r.GET("/publicnodes", func(c *gin.Context) {
		ip := c.DefaultQuery("ip", "")
		if len(ip) > 0 && net.ParseIP(ip) == nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"error": "please type right ip address!",
			})
			return
		}

		res, _ := storage.PublicNodeList(ip)
		c.HTML(http.StatusOK, "publicnodes.tmpl", gin.H{
			"peers": res,
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
