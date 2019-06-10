package main

import (
	"bufio"
	"os"
	"strings"
	"log"
)

//get all of the people a user if following
func getFriends(username string, lock *Lock) []string {
	//acquires a shared lock for friends.txt
	acquireLock("read", lock)
	
	friendsFile, err := os.Open(friendsFileName)
	
	if err != nil{
		log.Println("Could not open friends file properly.")
		log.Fatal(err)
	}
	
	friendsScanner := bufio.NewScanner(friendsFile)

	var friends []string

	//goes through friends file and if you're their friend, adds friend to array
	for friendsScanner.Scan() {

		line := friendsScanner.Text()
		lineInfo := strings.SplitAfter(line, ",")

		if len(lineInfo) >= 1 && lineInfo[0] != "" {
			currUsername := lineInfo[0][0:len(lineInfo[0])-1]

			if username == currUsername {
				friendUsername := lineInfo[1]
				friends = append(friends, friendUsername)
			}
		}

	}
	friendsFile.Close()
	releaseLock("read",lock)
	
	return friends
}

//if newsfeed, show your posts and your friend's posts
//if profile show only your posts
func retrievePosts(page, username string, lock *Lock) string {

	var users []string
	if page == "newsfeed" {
		users = getFriends(username, lock)
	} 
	
	//take all friends (if profile, no friends in array) and add your name as well
	users = append(users, username)

	//acquires a shared lock to read posts.txt
	acquireLock("read", lock)

	postFile, err := os.Open(postsFileName)
	
	if err != nil{
		log.Println("Could not open posts file properly.")
		log.Fatal(err)
	}
	
	postScanner := bufio.NewScanner(postFile)

	var userPosts []string

	//checks all posts, if the name of user is in users array (created right before this)
	for postScanner.Scan() {

		line := postScanner.Text()

		postInfo := strings.SplitAfter(line, ",")
		if len(postInfo) >= 1 && postInfo[0] != "" {
			lineUsername := postInfo[0][0:len(postInfo[0])-1]
		
			for i := 0; i < len(users); i++ {
				if users[i] == lineUsername {
					userPosts = append(userPosts, line)
				}
			}
		}
	}

	postFile.Close()
	releaseLock("read",lock)

	//will send back all messages as one huge string. will be formatted as follows:
	// user1: post1; user2: post2...userN: postN

	var allPosts string
	if len(userPosts) == 0 {
		allPosts = "No posts to display"

	} else {
		postArr := strings.SplitAfter(userPosts[0],",")
		userPost := postArr[0][0:len(postArr[0])-1]
		thePost := postArr[2]
		allPosts = userPost + ": " + thePost

		for i := 1; i < len(userPosts); i++ {
			allPosts = allPosts + ";"
			postArr = strings.SplitAfter(userPosts[i],",")
			userPost = postArr[0][0:len(postArr[0])-1]
			thePost = postArr[2]
			allPosts = allPosts + userPost + ": " + thePost
		}
	}

	return allPosts

}