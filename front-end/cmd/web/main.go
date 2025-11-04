package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// parse all page templates and partials once and set as Gin's template
	tmpl := template.Must(template.ParseGlob("./cmd/web/template/*.gohtml"))
	r.SetHTMLTemplate(tmpl)

	// routes
	r.GET("/", func(c *gin.Context) {
		render(c, "index.gohtml")
	})
	r.GET("/login", func(c *gin.Context) {
		render(c, "login.gohtml")
	})
	r.GET("/register", func(c *gin.Context) {
		render(c, "register.gohtml")
	})

	// run server
	_ = r.Run(":3000")
}

func render(c *gin.Context, t string) {
	partials := []string{
		"./cmd/web/template/base.layout.gohtml",
		"./cmd/web/template/header.partial.gohtml",
		"./cmd/web/template/footer.partial.gohtml",
	}

	page := fmt.Sprintf("./cmd/web/template/%s", t)

	files := []string{page}
	files = append(files, partials...)

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_, _ = c.Writer.Write([]byte(err.Error()))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(c.Writer, "base", nil); err != nil {
		c.Status(http.StatusInternalServerError)
		_, _ = c.Writer.Write([]byte(err.Error()))
		return
	}
}
