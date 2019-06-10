package main

import (
    "net/http"
    "fmt"
    "os"
    "bufio"
    "log"
    "net"
)

//shows remove account steps, esnure user wanted to click that, if so proceed
//requests BackEnd to remove all things having to do with this user
func fremoveAcct(w http.ResponseWriter, r *http.Request) {
	
	userN, err := r.Cookie("username")
	if err != nil {
		log.Println("Error retrieving cookie.")
	}

	if r.Method == "GET" {

		//render HTML which shows a yes or no button with making sure user wants
		//delete account
    	renderHtml("./removeAcct.html", w)

	} else {

		r.ParseForm()

		//if user confirms they want to remove their account
		if r.PostFormValue("Yes") == "Yes" {


			username := userN.Value

			//preparing message to send to request the BackEnd to remove account
			s := "removeAccount," + username

			//asks to connect to the server
			conn, err := net.DialTimeout("tcp", service, TIMEOUT)
			
			var newMasterID int

			for err != nil {
				fmt.Fprintf(os.Stderr, "could not connect", err.Error())
				log.Println("Server %i down",master.server_id)
				newMasterID = chooseNewMaster()
				conn, err = net.DialTimeout("tcp", service, TIMEOUT)
			}

			log.Println("After sending removeAccount request the master is %s", newMasterID)
			fmt.Fprintf(conn, "%s\n",s)

			//stores the BackEnd's reply 
			reply, err1 := bufio.NewReader(conn).ReadString('\n')
			if err1 != nil {
				log.Println("Remove Account reply could not be processed")
			}

			log.Println(reply)

			//after user is deleted, making sure to logout
			flogout(w,r)

		} else {

			http.Redirect(w,r,"/newsfeed",http.StatusPermanentRedirect)

		}

	}
}