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

type Contacts = []Contact

type Data struct {
    Contacts Contacts
}

type Page struct {
    Data Data
    Form FormData 
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newPage() Page {
    return Page{
        Data: newData(),
        Form: newFormData(),
    }
}

func newFormData() FormData {
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

func (d *Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if contact.Email == email {
            return true
        }
    } 
    return false
}

func newContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}

func newData() Data {
    return Data{
        Contacts: []Contact{
            newContact("Arnold", "arnold.lei88@gmail.com"),
            newContact("Jenny", "catscribbler22@gmail.com"),
        },
    }
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

    page := newPage()
    e.Renderer = newTemplate()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page) 
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email 
            formData.Errors["email"] = "Email already exists"

            c.Render(422, "form", formData)
        }
        contact := newContact(name, email)
        page.Data.Contacts = append(page.Data.Contacts, contact)
        c.Render(200, "form", newFormData())
        return c.Render(200, "oob-contact", contact) 
    })
    
    e.DELETE("/", func(c echo.Context) error {
        return c.Render(200, "index", page) 
    })

    e.Logger.Fatal(e.Start(":42069"))

}

