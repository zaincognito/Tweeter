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

//displays all of signup page, sends request to BackEnd to write new user
func fsignup(w http.ResponseWriter, r *http.Request) {
	
	//checks if cookie is empty, if it is not, illegal access to sign up
	c, err := r.Cookie("username")
	
	if err != nil {
		log.Println(err)
	}
	
	if c == nil || c.Value == "" {
		if r.Method == "GET" {

	    	//Render the page
		    renderHtml("./signup.html", w)

		//this else is processing the sign-up form input
		} else {

			r.ParseForm()

			//Get the information from the form so we can store it
			username := strings.Join(r.Form["user"], " ")
			password := strings.Join(r.Form["pass"], " ")
			name := strings.Join(r.Form["name"], " ")

			//if any of these fields are empty, try again
			if(username == "" || password == "" || name == "") {
				http.Redirect(w,r,"/signup",302)
			}


			//prepare message to send to request BackEnd to sign this user up
			s := "signup," + username + "," + password + "," + name

			//ask to connect to service
			conn, err := net.DialTimeout("tcp", service,TIMEOUT)
			var newMasterID int
			for err != nil {
				fmt.Fprintf(os.Stderr, "could not connect", err.Error())
				log.Println("TRIED CONNECTING BUT SERVER %i DOWN",master.server_id)
				newMasterID = chooseNewMaster()
				conn, err = net.DialTimeout("tcp", service, TIMEOUT)
			}
			
			log.Println("AFTER SIGNUP REQUEST THE MASTER IS %i", newMasterID)
			fmt.Fprintf(conn, "%s\n",s)
			
			//read reply given by server
			reply, err1 := bufio.NewReader(conn).ReadString('\n')
			
			if err1 != nil {
				log.Println("Signup reply could not be processed")
			}

			//BackEnd successfully created the user, create a cookie and redirect
			//to user's homepage
			if reply == "Yes\n" {

				setCookie(w,username, 60)
				log.Println("Finished setting cookie")
				http.Redirect(w,r,"/newsfeed",http.StatusPermanentRedirect)
				
			} else {
				http.Redirect(w,r,"/signup",302)
			}
		}
	} else {
		http.Redirect(w,r,"/newsfeed",302)
	}
}