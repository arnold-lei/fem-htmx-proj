package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

type Count struct {
    Count int
}
type Contact struct {
    Name string
    Email string
}
func (c *Contact) newContact(name string, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}

type Contacts = []Contact

type Data struct {
    Contacts Contacts
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

func main()  {

    e := echo.New() 
    e.Use(middleware.Logger())

    count := Count { Count: 0 }
    e.Renderer = newTemplate()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", count) 
    })
    e.POST("/contacts", func(c echo.Context) error {
        return c.Render(200, "index", count) 
    })
    e.Logger.Fatal(e.Start(":42069"))

}

