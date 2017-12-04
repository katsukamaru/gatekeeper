package main

import (
	"fmt"
	"gatekeeper/keymanage"
	"html/template"
	"net/http"
)

type Server struct {
	Name  string
	Users []User
}

type User struct {
	Name string
	Sudo bool
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	user1 := User{"user1", true}
	users := []User{user1}
	server := Server{"server1", users}
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, server)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, World")

}
func userAdd_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form.Get("userlist"))
	s := r.Form.Get("userlist")
	keymanage.UserAdd(s)
	fmt.Fprintf(w, "userAdd_handler")
}

func userDel_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := r.Form.Get("userlist")
	keymanage.UserDel(s)
	fmt.Fprintf(w, "userDel_handler")
}

func authAdd_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := r.Form.Get("username")
	keymanage.AuthAdd(s)
	fmt.Fprintf(w, "{\"response\" : \"AuthAdd_handler\"}")
}

func authDel_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "authDel_handler")
}

func initPass_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "initPass_handler")
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/useradd", userAdd_handler)
	http.HandleFunc("/userdel", userDel_handler)
	http.HandleFunc("/authadd", authAdd_handler)
	http.HandleFunc("/authdel", authDel_handler)
	http.HandleFunc("/initpass", initPass_handler)
	http.ListenAndServe(":8080", nil)
}
