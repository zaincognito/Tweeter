package main

import (
	"net"
	"net/http"
	"strings"
	"fmt"
	"bufio"
	"os"
	"log"
)

//displays signin page, which requests BackEnd once the signin is clicked to
//validate user and if it works, create a cookie and send to user's homepage
func fsignin(w http.ResponseWriter, r *http.Request) {
	
	c, err := r.Cookie("username")
	
	//if cookie already exists, that means the user went back to signin with out
	//logging out
	if err != nil {
		log.Println(err)
	}
	
	if c == nil || c.Value == "" {
		
		if r.Method == "GET" {
	    	
	    	//Render the HTML page which has a form to signin
	    	renderHtml("./signin.html", w)

		} else {
			
			r.ParseForm()

			//get the username and password from the Form
			username := strings.Join(r.Form["user"], " ")
			password := strings.Join(r.Form["pass"], " ")


			//preparing message to send to request BackEnd to check if user is valid
			s := "signin," + username + "," + password

			//asks to connect to server
			conn, err := net.DialTimeout("tcp", service, TIMEOUT)
			
			var newMasterID int

			for err != nil {
				fmt.Fprintf(os.Stderr, "could not connect", err.Error())
				log.Println("Server %i down",master.server_id)
				newMasterID = chooseNewMaster()
				conn, err = net.DialTimeout("tcp", service, TIMEOUT)
			}

			log.Println("After sending signin request the master is %s", newMasterID)
			fmt.Fprintf(conn, "%s\n",s)

			//Reply from BackEnd stored. If "yes", we can proceed to signing in
			//if not, try again
			reply, err1 := bufio.NewReader(conn).ReadString('\n')
			
			if err1 != nil {
				log.Println("Signin reply could not be processed")
			}
			
			if reply == "Yes\n" {
				
				setCookie(w,username,5)
				http.Redirect(w,r,"/newsfeed",302)
			} else {
				http.Redirect(w,r,"/signin",302)
			}
		}
	} else {
		http.Redirect(w,r,"/newsfeed",302)
	}
}