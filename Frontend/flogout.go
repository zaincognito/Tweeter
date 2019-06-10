package main

import (
	"net/http"
)

//logs user out->deletes cookie by aging it and making it invalid, redirects to
//signin page
func flogout(w http.ResponseWriter, r *http.Request) {
	
	c, _ := r.Cookie("username")
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	http.SetCookie(w,c)
	http.Redirect(w,r,"/signin",302)

}