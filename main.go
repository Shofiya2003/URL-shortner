package main

import (
	"fmt"
	"hash/crc32"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Url struct {
	LongUrl  string
	Key      string
	ShortUrl string
}

func pingServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "we are logically blessed!"})
}

func postURL(c *gin.Context) {
	var newUrl Url

	err := c.BindJSON(&newUrl)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	if newUrl.LongUrl == "" {
		fmt.Println("missing long url")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": "url missing"})
		return
	}

	longUrl := newUrl.LongUrl
	crc32q := crc32.MakeTable(0xD5828281)
	hash := crc32.Checksum([]byte(longUrl), crc32q)

	// add logic to handle duplicate key from a different long url
	hashString := fmt.Sprintf("%08x", hash)
	// check if the long url already exists
	fetchedUrl, err := FetchUrl(hashString)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"err": "could not fetch url"})
		return
	}

	if fetchedUrl.LongUrl != newUrl.LongUrl {
		newUrl.Key = hashString

		newUrl.ShortUrl = fmt.Sprintf("http://localhost:8000/%x", hash)

		err = StoreUrl(newUrl)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}
	} else {
		fmt.Println("Long URL already exists")
		newUrl.Key = fetchedUrl.Key
		newUrl.ShortUrl = fetchedUrl.ShortUrl
	}

	c.IndentedJSON(http.StatusOK, gin.H{"key": newUrl.Key, "long_url": newUrl.LongUrl, "short_url": newUrl.ShortUrl})
}

func main() {
	router := gin.Default()

	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}

	router.POST("/", postURL)
	router.GET("/ping", pingServer)
	router.Run("localhost:8000")

}
