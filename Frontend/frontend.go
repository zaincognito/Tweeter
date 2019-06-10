package main

import (
	"net/http"
	"html/template"
	"time"
)

var TIMEOUT time.Duration = time.Second
var service string
var repManagers ServerQueue
var master BEServer

//takes an HTML file name and executes it to write to the response writer
//and display
func renderHtml(htmlFile string, w http.ResponseWriter){
	temp, _ := template.ParseFiles(htmlFile)
	temp.Execute(w,nil)
}

//creates a cookie when a user will have a session
func setCookie(w http.ResponseWriter, username string, minutes time.Duration) {
	expiration := time.Now().Add(minutes*time.Minute)
	cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func main() {

	repManagers = populateServers(repManagers)
	master = repManagers.Front()
	service = master.port

	//goes to each respective function to handle displaying and requesting
	//BackEnd to give information
	http.HandleFunc("/", fsignin)
	http.HandleFunc("/signin", fsignin)
	http.HandleFunc("/signup", fsignup)
	http.HandleFunc("/newsfeed", fnewsfeed)
	http.HandleFunc("/logout", flogout)
	http.HandleFunc("/profile", fprofile)
	http.HandleFunc("/removeAcct", fremoveAcct)

	//Listening on port in order to host the FrontEnd
	http.ListenAndServe(":7050", nil)
}