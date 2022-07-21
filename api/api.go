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
	r.GET("/query/hash", func(c *gin.Context) {
		hash := c.Query("hash")
		hash = strings.ToUpper(hash)
		if hash == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No hash provided",
			})
			return
		} else if len(hash) != 40 && len(hash) != 32 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid hash provided",
			})
			return
		}
		match := false
		for key, blooms := range buckets.BloomFilters {
			if contains, _ := blooms.MightContain([]byte(hash)); contains {
				f, _ := os.Open("buckets/" + key + ".bkt")
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {

					if strings.Contains(scanner.Text(), hash) {
						match = true
						delimitedText := strings.Split(strings.ReplaceAll(scanner.Text(), "\"", ""), ",")
						tempNsrlDataPoint := db.NSRLDataPoint{
							SHA1:         delimitedText[0],
							MD5:          delimitedText[1],
							CRC32:        delimitedText[2],
							FileName:     delimitedText[3],
							FileSize:     delimitedText[4],
							ProductCode:  delimitedText[5],
							OpSystemCode: delimitedText[6],
							SpecialCode:  delimitedText[7],
						}
						c.JSON(http.StatusOK, gin.H{
							"message": tempNsrlDataPoint,
						})
					}
				}
			}
		}
		if !match {
			c.JSON(http.StatusOK, gin.H{
				"message": "No matches found",
			})
		}
	})
	r.GET("/query/file", func(c *gin.Context) {
		file := c.Query("file")
		if file == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file provided",
			})
			return
		}
		match := false
		returnData := []db.NSRLDataPoint{}
		for key, blooms := range buckets.BloomFilters {
			if contains, _ := blooms.MightContain([]byte(file)); contains {
				f, _ := os.Open("buckets/" + key + ".bkt")
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					delimitedText := strings.Split(strings.ReplaceAll(scanner.Text(), "\"", ""), ",")
					if len(delimitedText) != 8 {
						continue
					}
					if delimitedText[3] == file {
						match = true
						tempNsrlDataPoint := db.NSRLDataPoint{
							SHA1:         delimitedText[0],
							MD5:          delimitedText[1],
							CRC32:        delimitedText[2],
							FileName:     delimitedText[3],
							FileSize:     delimitedText[4],
							ProductCode:  delimitedText[5],
							OpSystemCode: delimitedText[6],
							SpecialCode:  delimitedText[7],
						}
						returnData = append(returnData, tempNsrlDataPoint)
					}

				}
				c.JSON(http.StatusOK, gin.H{
					"message": returnData,
				})
			}
		}
		if !match {
			c.JSON(http.StatusOK, gin.H{
				"message": "No matches found",
			})
		}
	})

	r.Run() // listen and serve on port 8080
}
