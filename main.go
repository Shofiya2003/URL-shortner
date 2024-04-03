package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type Url struct {
	LongUrl  string
	Key      string
	ShortUrl string
}

type ResultUrl struct {
	Index     int    `json:"id"`
	Long_url  string `json:"long_url"`
	Short_url string `json:"short_url"`
}

func PingServer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "url shortner apis!"})
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
	fmt.Println(hash)
	// add logic to handle duplicate key from a different long url
	hashString := fmt.Sprintf("%08x", hash)
	// check if the long url already exists
	fetchedUrl, err := FetchUrl(hashString)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"err": "could not fetch url"})
		return
	}

	fmt.Println(hashString)
	fmt.Println(hash)
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

func redirect(c *gin.Context) {
	key := c.Param("key")

	fetchedUrl, err := FetchUrl(key)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"err": "could not fetch url"})
		return
	}

	if fetchedUrl.LongUrl == "" {
		fmt.Println("URL not found")
		c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "URL not found"})
		return
	}
	c.Header("Location", fetchedUrl.LongUrl)
	c.Status(302)

}

func deleteUrl(c *gin.Context) {
	key := c.Param("key")

	err := DeleteUrl(key)

	if err != nil {
		fmt.Printf("error in deleting key %s ", key)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func listAllUrl(c *gin.Context) {
	pageQuery := c.Query("page")

	if pageQuery == "" {
		pageQuery = "1"
	}
	page, err := strconv.Atoi(pageQuery)

	if err != nil {
		log.Println("error in getting page query")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "incorrect page query"})
		return
	}
	urls, err := FetchAllUrl(page)
	if err != nil {
		log.Println("error in fetching all urls")
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error in fetching urls"})
		return
	}

	result := make([]ResultUrl, len(urls))
	start := (page - 1) * 2
	for i, url := range urls {
		result[i] = ResultUrl{start + i + 1, url.LongUrl, url.ShortUrl}
	}

	c.IndentedJSON(200, result)
}

func main() {
	router := gin.Default()

	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}

	router.GET("/ping", PingServer)
	router.POST("/", postURL)
	router.GET("/:key", redirect)
	router.POST("/all", listAllUrl)
	router.DELETE("/:key", deleteUrl)
	router.Run("localhost:8000")

}
