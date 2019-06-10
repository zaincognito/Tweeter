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

//display all of a user's profile page, requesting BackEnd for the friends
//they have, the posts, and if they want to write a post
func fprofile(w http.ResponseWriter, r *http.Request) {
	
	//check if cookie exists already, if not then illegally entered
	userN, err := r.Cookie("username")
	
	if(err == nil) {
	    
	    username := userN.Value
	    
	    //render HTML to display navigation bar on top
	    renderHtml("./navbar.html", w)


	    //preparing message to send to request my Name
		s := "usersName," + username

		//ask to connect to server
		conn, err := net.DialTimeout("tcp", service, TIMEOUT)
		
		var newMasterID int

		for err != nil {
			fmt.Fprintf(os.Stderr, "could not connect", err.Error())
			log.Println("Server %i down",master.server_id)
			newMasterID = chooseNewMaster()
			conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		}

		log.Println("After sending usersName request the master is %s", newMasterID)
		fmt.Fprintf(conn, "%s\n",s)

		//get reply back from BackEnd with our Name
	 	name, err := bufio.NewReader(conn).ReadString('\n')
		
		if err != nil {
			log.Println("UsersName reply could not be processed")
		}

		//display Name
		fmt.Fprintf(w, "<h1>%s</h1>", name)

		//preparing message to send to request BackEnd to send all of our
		//friend's usernames
		s = "getFriends," + username

		//ask to connect to server
		conn, err = net.DialTimeout("tcp", service, TIMEOUT)

		for err != nil {
			fmt.Fprintf(os.Stderr, "could not connect", err.Error())
			log.Println("Server %i down",master.server_id)
			newMasterID = chooseNewMaster()
			conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		}

		log.Println("After sending getFriends request the master is %s", newMasterID)
		fmt.Fprintf(conn, "%s\n",s)

		//get reply back with all of my friends in one string
		friends, err := bufio.NewReader(conn).ReadString('\n')
		
		if err != nil {
			log.Println("getFriends reply could not be processed")
		}
		
		if friends == "No friends\n" {
			fmt.Fprintf(w,friends)

		//display friends
		} else {
			fmt.Fprintf(w,"<h3>Friends:</h3>")
			fmt.Fprintf(w,friends)
		}
		fmt.Fprintf(w,"<h3> </h3>")

		//form to search for a username to remove from friend's list (unfollow)
		fmt.Fprintf(w, "<form method = \"post\"> Remove friends by Username:   "+ 
		"<input type=\"text\" name=\"removed\"/>"+ 
		"<input type=\"submit\" value=\"Remove\"/></br></br>"+ 
		"</form>")

		//if unfollow clicked
		if r.Method == "POST" {
			
			r.ParseForm()

			//Get username that user wants to unfollow
			removeFriend := r.PostFormValue("removed")

			//if the search box was not empty
			if removeFriend != "" {
				

				//preparing message to send to request BackEnd to delete friend
				s = "remFriend," + username + "," + removeFriend

				//ask to connect to server
				conn, err := net.DialTimeout("tcp", service,TIMEOUT)
				
				var newMasterID int
				
				for err != nil {
					fmt.Fprintf(os.Stderr, "could not connect", err.Error())
					log.Println("Server %i down",master.server_id)
					newMasterID = chooseNewMaster()
					conn, err = net.DialTimeout("tcp", service, TIMEOUT)
				}
				
				log.Println("After sending removeFriend request the master is %s", newMasterID)
				fmt.Fprintf(conn, "%s\n",s)

				//get reply back with message whether unfollow worked or not
				remFriendMessage, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					log.Println("RemFriend reply could not be processed")
				}

				//display unfollow message
				fmt.Fprintf(w,  "%s<br>",remFriendMessage)
			}

		}

		//renter HTML which has a way for the user to post
		renderHtml("./profile.html", w)

		//if post was clicked
		if r.Method == "POST" {
			
			r.ParseForm()

			thePost := strings.Join(r.Form["aPost"], " ")

			//if the tweet was not empty
			if thePost != "" {

				//preparing message to send to request the BackEnd to write post
				s = "writePost," + username + "," + thePost

				//ask to connect to server
				conn, err := net.DialTimeout("tcp", service,TIMEOUT)
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

		//preparing message to send to request BackEnd to send all user's posts
		s = "retrievePosts," + "profile," + username

		//ask to connect to server
		conn, err = net.DialTimeout("tcp", service, TIMEOUT)

		for err != nil {
			fmt.Fprintf(os.Stderr, "could not connect", err.Error())
			log.Println("Server %i down",master.server_id)
			newMasterID = chooseNewMaster()
			conn, err = net.DialTimeout("tcp", service, TIMEOUT)
		}

		log.Println("After sending retrievePosts request the master is %s", newMasterID)
		fmt.Fprintf(conn, "%s\n",s)

		//get reply back with all user posts
		allPosts, err := bufio.NewReader(conn).ReadString('\n')
		
		if err != nil {
			log.Println("retrievePosts reply could not be processed")
		}
		
		if allPosts == "No posts to display\n" {
			log.Println(allPosts)

		//display all of the posts in a nice fashion
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