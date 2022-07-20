package api

import (
	"bufio"
	"net/http"
	"nsrlFerret/db"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func StartApi(buckets db.BucketCollection) {
	startTime := time.Now()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to NSRL Ferret",
		})
	})
	r.GET("/stats/buckets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"bucket_count": len(buckets.Uuids),
		})
	})
	r.GET("/stats/uptime", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uptime": time.Since(startTime).String(),
		})
	})

	r.GET("/query", func(c *gin.Context) {
		hash := c.Query("hash")
		for key, blooms := range buckets.BloomFilters {
			if contains, _ := blooms.MightContain([]byte(hash)); contains {
				f, _ := os.Open("buckets/" + key + ".bkt")
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), hash) {
						c.JSON(http.StatusOK, gin.H{
							"message": scanner.Text(),
						})
					}
				}
			}
		}
	})

	r.Run() // listen and serve on port 8080
}
