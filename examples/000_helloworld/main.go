package main

import (
	"net/http"

	"github.com/lixinio/kelly"
)

func main() {
	router := kelly.New(nil)

	router.GET("/", func(c *kelly.Context) {
		c.WriteIndentedJSON(http.StatusOK, kelly.H{
			"code": "0",
		})
	})

	router.Run(":9999")
}
