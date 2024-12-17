package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func translateWord(word string) (string, error) {
	baseURL := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=tr&tl=en&dt=t&q="
	resp, err := http.Get(baseURL + word)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.Contains(string(body), `<title>Error 400 (Bad Request)`) {
		return "", fmt.Errorf("invalid request")
	}

	var result []interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if len(result) > 0 {
		translatedText := ""
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translated := range slice.([]interface{}) {
				translatedText = fmt.Sprintf("%v", translated)
				break
			}
		}
		return translatedText, nil
	}
	return "", fmt.Errorf("no translation found")
}

func main() {
	r := gin.Default()

	r.GET("/translate/:word", func(c *gin.Context) {
		word := c.Param("word")

		translatedWord, err := translateWord(word)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": translatedWord,
		})
	})
	 r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
 r.GET("/hello", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "world!",
    })
  })
	
	r.Run(":80")
}

