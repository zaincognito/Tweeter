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

//displays all of a user's newsfeed by making requests to BackEnd to send
//user's info and posts
func fnewsfeed(w http.ResponseWriter, r *http.Request) {
	
	//Render the page
	userN, err := r.Cookie("username")

	//checks if there is an error with cookie (happens when illegally in)
	if(err == nil) {
	    
	    username := userN.Value

	    //render HTML for the navigation bar that is always on profile
	    renderHtml("./navbar.html", w)


		//preparing message to send to request BackEnd to give my Name
		s := "usersName," + username

		//asks to connect to server
		conn, err := net.DialTimeout("tcp", service,TIMEOUT)
		
		var newMasterID int

		for err != nil {
			fmt.Fprintf(os.Stderr, "could not connect", err.Error())
			log.Println("Server %i down",master.server_id)
			newMasterID = chooseNewMaster()
			conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		}

		log.Println("After sending usersName request the master is %s", newMasterID)
		fmt.Fprintf(conn, "%s\n",s)

		//reads in name from reply from the BackEnd
		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("UsersName reply could not be processed")
		}

		//Displays name
		fmt.Fprintf(w, "<h1>Welcome %s!</h1>", name)

		//form to search for a username to add to their friend's list (follow)
		fmt.Fprintf(w, "<form method = \"post\"> Search for Friends by Username:   "+ 
			"<input type=\"text\" name=\"search\"/>"+ 
			"<input type=\"submit\" value=\"Follow\"/></br></br>"+ 
			"</form>")

		//process the search of a username to add as friend
		if r.Method == "POST"{
			r.ParseForm()

			//Get what the user inputted as "friend's" username
			prospectFriendUsername := r.PostFormValue("search")

			//if the search box isn't empty
			if prospectFriendUsername != "" {
				

				//preparing message to send to request a friend to be added
				s = "addFriend," + username + "," + prospectFriendUsername

				//asks to connect to server
				conn, err = net.DialTimeout("tcp", service,TIMEOUT)
				
				var newMasterID int

				for err != nil {
					fmt.Fprintf(os.Stderr, "could not connect", err.Error())
					log.Println("Server %i down",master.server_id)
					newMasterID = chooseNewMaster()
					conn, err = net.DialTimeout("tcp", service, TIMEOUT)
				}

				log.Println("After sending addFriend request the master is %s", newMasterID)
				fmt.Fprintf(conn, "%s\n",s)

				//reads in reply from BackEnd to confirm if adding frienf worked
				addFriendMessage, err := bufio.NewReader(conn).ReadString('\n')
				
				if err != nil {
					log.Println("AddFriend reply could not be processed")
				}

				//display message about adding friend
				fmt.Fprintf(w,  "%s<br>",addFriendMessage)
				fmt.Fprintf(w,  "<h3> </h3>")
			}
		}

		//render HTML page for newdfeed->writing a post
		renderHtml("./newsfeed.html", w)

		//If the user clicked post
		if r.Method == "POST" {
			
			r.ParseForm()

			thePost := strings.Join(r.Form["aPost"], " ")

			//if it is not an empty string, request BackEnd to write post
			if thePost != "" {

				//preparing message to send to request BackEnd to write post
				s = "writePost," + username + "," + thePost

				//asks to connect to the server
				conn, err = net.DialTimeout("tcp", service,TIMEOUT)
				
				var newMasterID int
				for err != nil {
					fmt.Fprintf(os.Stderr, "could not connect", err.Error())
					log.Println("Server %i down",master.server_id)
					newMasterID = chooseNewMaster()
					conn, err = net.DialTimeout("tcp", service, TIMEOUT)
				}
				
				log.Println("After sending writePost request the master is %s", newMasterID)
				fmt.Fprintf(conn, "%s\n",s)
			}

		}

		//preparing message to send to request BackEnd for the posts
		s = "retrievePosts," + "newsfeed," + username

		//ask to connect
		conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		
		for err != nil {
			fmt.Fprintf(os.Stderr, "could not connect", err.Error())
			log.Println("Server %i down",master.server_id)
			newMasterID = chooseNewMaster()
			conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		}

		log.Println("After sending newsfeed posts request the master is %s", newMasterID)
		fmt.Fprintf(conn, "%s\n",s)

		//read in reply from BackEnd with the posts needed to display
		allPosts, err := bufio.NewReader(conn).ReadString('\n')
		
		if err != nil {
			log.Println("retrievePosts reply could not be processed")
		}
		
		if allPosts == "No posts to display\n" {
			log.Println(allPosts)

		//display all posts in a nice fashion
		} else {
			postsArr := strings.SplitAfter(allPosts,";")
			for i:= 0; i < len(postsArr); i++ {
				curPost := postsArr[i][0:len(postsArr[i])-1]
				fmt.Fprintf(w,"<h2>%s</h2>",curPost)
			}
		}
		
	} else {
		http.Redirect(w,r,"/signin",http.StatusPermanentRedirect)
	}
}