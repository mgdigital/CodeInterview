package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "assets.db"

type Asset struct {
	ID        int
	Host      string
	Comment   string
	Owner     string
	IPs       []IP
	Ports     []Port
	Signature string
}

type IP struct {
	Address   string
	Signature string
}

type Port struct {
	Port      int
	Signature string
}

func generateSignature(asset Asset) Asset {
	data := asset.Host + asset.Comment + asset.Owner
	hash := sha256.New()
	hash.Write([]byte(data))
	signature := hex.EncodeToString(hash.Sum(nil))

	newAsset := asset
	newAsset.Signature = signature

	for i, ip := range newAsset.IPs {
		ipData := ip.Address
		hash := sha256.New()
		hash.Write([]byte(ipData))
		ip.Signature = hex.EncodeToString(hash.Sum(nil))

		newAsset.IPs[i] = ip
	}

	for i, port := range newAsset.Ports {
		portData := fmt.Sprintf("%d", port.Port)
		hash := sha256.New()
		hash.Write([]byte(portData))
		port.Signature = hex.EncodeToString(hash.Sum(nil))

		newAsset.Ports[i] = port
	}

	return newAsset
}

func main() {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/assets", func(c *gin.Context) {
		assetID := c.Query("id")
		var rows *sql.Rows
		if assetID != "" {
			rows, err = db.Query("SELECT id, host, comment, owner FROM assets WHERE id = ?", assetID)
		} else {
			rows, err = db.Query("SELECT id, host, comment, owner FROM assets")
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var id int
			var host, comment, owner string
			if err := rows.Scan(&id, &host, &comment, &owner); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var ips []IP
			ipRows, err := db.Query("SELECT address FROM ips WHERE asset_id = ?", id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer ipRows.Close()
			for ipRows.Next() {
				var ipAddress string
				if err := ipRows.Scan(&ipAddress); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				ips = append(ips, IP{Address: ipAddress})
			}

			var ports []Port
			portRows, err := db.Query("SELECT port FROM ports WHERE asset_id = ?", id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer portRows.Close()
			for portRows.Next() {
				var portNum int
				if err := portRows.Scan(&portNum); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				ports = append(ports, Port{Port: portNum})
			}

			asset := Asset{
				ID:      id,
				Host:    host,
				Comment: comment,
				Owner:   owner,
				IPs:     ips,
				Ports:   ports,
			}

			processedAsset := generateSignature(asset)

			assets = append(assets, processedAsset)
		}

		c.JSON(http.StatusOK, assets)
	})

	router.Run(":8080")
}
