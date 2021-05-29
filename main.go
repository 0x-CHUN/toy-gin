package main

import (
	"go-web/spring"
	"log"
	"net/http"
	"time"
)

func main() {
	r := spring.New()
	r.Use(spring.Logger())
	r.Use(spring.Recovery())

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/index", func(c *spring.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/panic", func(c *spring.Context) {
		names := []string{"test"}
		c.String(http.StatusOK, names[100])
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *spring.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

	}
	v2 := r.Group("/v2")
	v2.Use(Logger2())
	{
		v2.GET("/hello/:name", func(c *spring.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *spring.Context) {
			c.JSON(http.StatusOK, spring.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}

func Logger2() spring.HandlerFunc {
	return func(c *spring.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
