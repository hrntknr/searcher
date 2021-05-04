package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func newController(config *config, service Service) (*controller, error) {
	router := gin.New()

	router.POST("/regist", func(c *gin.Context) {
		var body RegistBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := service.Regist(body.Uri, body.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, nil)
	})

	router.GET("/search", func(c *gin.Context) {
		var count uint
		if c.Query("count") == "" {
			count = 10
		} else {
			_count, err := strconv.ParseUint(c.Query("count"), 10, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			count = uint(_count)
		}
		var offset uint
		if c.Query("offset") == "" {
			offset = 0
		} else {
			_offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			offset = uint(_offset)
		}
		result, err := service.Search(c.Query("k"), offset, count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	})

	return &controller{
		router: router,
		config: config,
	}, nil
}

type controller struct {
	router *gin.Engine
	config *config
}

func (c *controller) start() error {
	err := c.router.Run(c.config.Listen)
	if err != nil {
		return err
	}
	return nil
}

type RegistBody struct {
	Uri  string `json:"uri"`
	Body string `json:"body"`
}
