package main

import (
	"html/template"
	"net/http"
	"log"
	"fmt"
)

var tpl *template.Template

func init() {
	tpl, _ = template.ParseGlob("templates/*.html")
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/imgs/", fs)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie("session-id")
	if err != nil {
		cookie = newVisitor()
		http.SetCookie(res, cookie)
		fmt.Println("no cookie newvisitor ran") // DEBUGGING
	}

	if tampered(cookie.Value) {
		fmt.Println("INSIDE IF TAMPERED COOKIE") // DEBUGGING
		cookie = newVisitor()
		http.SetCookie(res, cookie)
		fmt.Println("tampered cookie newvisitor ran") // DEBUGGING
	}

	if req.Method == "POST" {
		src, hdr, err := req.FormFile("data")
		if err != nil {
			log.Println("error uploading photo: ", err)
			// TODO: create error page to show user
		}
		cookie = uploadPhoto(src, hdr, cookie)
		http.SetCookie(res, cookie)
		fmt.Println("upload photo ran") // DEBUGGING
	}

	m := Model(cookie)
	tpl.ExecuteTemplate(res, "index.html", m)
}