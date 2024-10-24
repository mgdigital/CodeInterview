package main

import (
	"errors"
	"fmt"
	"interview/assets"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "assets.db"

func main() {
	lookup, err := assets.NewLookup(dbFileName)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Cache-Control"},
		AllowCredentials: true,
	}))

	router.GET("/assets", func(c *gin.Context) {
		params, err := getParams(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		result, err := lookup.Lookup(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.Header("Cache-Control", "public, max-age=600")
		c.JSON(http.StatusOK, result)
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func getParams(c *gin.Context) (assets.Params, error) {
	limit := 10
	var err error
	if strLimit := c.Query("limit"); strLimit != "" {
		if limit, err = strconv.Atoi(strLimit); err != nil {
			return assets.Params{}, fmt.Errorf("invalid limit: %w", err)
		}
		if limit > 1000 {
			return assets.Params{}, errors.New("limit must be <= 1000")
		}
	}
	offset := 0
	if strOffset := c.Query("offset"); strOffset != "" {
		if offset, err = strconv.Atoi(strOffset); err != nil {
			return assets.Params{}, fmt.Errorf("invalid offset: %w", err)
		}
	}
	order := assets.SortOrderHost
	if strOrder := c.Query("order"); strOrder != "" {
		if parsed, ok := assets.ParseSortOrder(strOrder); !ok {
			return assets.Params{}, fmt.Errorf("invalid order: '%s'", strOrder)
		} else {
			order = parsed
		}
	}
	desc := c.Query("desc") == "true"
	return assets.Params{
		AssetID:   c.Query("id"),
		Query:     c.Query("q"),
		Limit:     limit,
		Offset:    offset,
		SortOrder: order,
		SortDesc:  desc,
	}, nil
}
