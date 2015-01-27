package main

import (
    "html/template"
    "net/http"
    "fmt"
    "io/ioutil"
)

const BASE_TEMPLATE = "templates/base.html"

var templates = map[string]*template.Template{
    "templates/time.html": nil,
    "templates/hello.html": nil,
}

func init() {
    for key, _ := range templates {
        templates[key] = template.Must(template.ParseFiles(
            "templates/base.html",
            key,
        ))
        fmt.Println("render", key)
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    // funcMap := template.FuncMap{
    //     "safehtml": func(text string) template.HTML { return template.HTML(text) },
    // }
    

    
}

func renderBaseTemplate(res http.ResponseWriter, templateLoc string, data interface{}) {
    tmpl, ok := templates[templateLoc]
    if ok {
        fmt.Println("ok!")
    }

    err := tmpl.ExecuteTemplate(res, "base", data)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
    }
}

func serveStaticFile(filename string) func(http.ResponseWriter, *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        buf, _ := ioutil.ReadFile(filename)
        res.Write(buf)
    }
}

func main() {
    fmt.Println("started!")
    http.HandleFunc("/time/", func(res http.ResponseWriter, req *http.Request) {
        data := struct {
            Time string
        }{
            Time: "12332434:3333",
        }
        fmt.Println("render /time/")
        renderBaseTemplate(res, "templates/time.html", data)
    }) 
    http.HandleFunc("/css/style.css", serveStaticFile("templates/style.css"))

    http.ListenAndServe(":8080", nil)
}
