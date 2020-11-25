package cloudscan

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var currentURL string

type RequestURL struct {
	URL string `json:"url"`
}

type Response struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

func add(c *gin.Context) {
	var urlJSON RequestURL
	err := c.BindJSON(&urlJSON)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	url := urlJSON.URL
	if url == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "no url present",
		})
		return
	}
	currentURL = url
	MessageInputChan <- currentURL
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func raw(c *gin.Context) {
	c.String(http.StatusOK, currentURL)
}

func redirect(c *gin.Context) {
	c.Redirect(302, currentURL)
}

func urlinfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"url": currentURL,
	})
}
