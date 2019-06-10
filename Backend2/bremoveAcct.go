package main

import (
	"bufio"
	"os"
	"strings"
	"log"
)

//removes any relationships where I am a user's friend, or I have a friend
func removeMeFriendFile(user string) {

	friendsFile, err := os.Open(friendsFileName)
	if err != nil{
		log.Println("Could not open friends file properly.")
		log.Fatal(err)
	}

	friendsScanner := bufio.NewScanner(friendsFile)

	var friends []string

	//goes through each line in friends.txt, if I am not in a line
	//store into an array
	for friendsScanner.Scan() {
		curUserFriend := friendsScanner.Text()
		curUserFriendArr := strings.SplitAfter(curUserFriend,",")
		if len(curUserFriendArr) >= 1 && string(curUserFriend[0]) != "" {
			curUser := curUserFriendArr[0][0:len(curUserFriendArr[0])-1]
			curFriend := curUserFriendArr[1]
			if (user != curUser) && (user != curFriend) {
				log.Println(curUser)
				friends = append(friends, curUserFriend)
			}
		}
	}

	//delete the whole file
	err = os.Remove(friendsFileName)
	if err != nil {
		log.Println("Could not delete friends file properly.")
	}

	//create friends file again
	newFriendsFile, err := os.Create(friendsFileName)
	if err != nil {
		log.Println("Could not create new friends file properly.")
	}

	newFriendsFile.Close()

	newFriendsFile, err = os.OpenFile(friendsFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}

	defer newFriendsFile.Close()

	//put back contents into file, while does not have me in it anymore
	putInFile := ""
	for i:=0; i < len(friends); i++ {
		putInFile = friends[i] + "\n"
		if _, err = newFriendsFile.WriteString(putInFile); err != nil {
		    panic(err)
		}
	}
}

//removes all of my posts
func removeMyPosts(user string) {

	postsFile, err := os.Open(postsFileName)
	if err != nil{
		log.Println("Could not open posts file properly.")
		log.Fatal(err)
	}

	postsScanner := bufio.NewScanner(postsFile)

	var posts []string

	//goes through posts.txt, stores any post that is not mine into array
	for postsScanner.Scan() {
		curPost := postsScanner.Text()
		postArr := strings.SplitAfter(curPost,",")
		if len(postArr) >= 1 && string(postArr[0]) != "" {
			curUser := postArr[0][0:len(postArr[0])-1]
			if user != curUser {
				posts = append(posts, curPost)
			}
		}
	}

	//delete entire posts.txt file
	err = os.Remove(postsFileName)
	if err != nil {
		log.Println("Could not delete posts file properly.")
	}

	//create posts.txt file again
	newPostsFile, err := os.Create(postsFileName)
	if err != nil {
		log.Println("Could not create new posts file properly.")
	}

	newPostsFile.Close()

	newPostsFile, err = os.OpenFile(postsFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}

	defer newPostsFile.Close()

	//put back contents of file, which does not have my posts in it anymore
	putInFile := ""
	for i:=0; i < len(posts); i++ {
		putInFile = posts[i] + "\n"
		if _, err = newPostsFile.WriteString(putInFile); err != nil {
		    panic(err)
		}
	}
}

//removes my account information
func removeMeUserFile(user string) {

	usersFile, err := os.Open(usersFileName)
	if err != nil{
		log.Println("Could not open users file properly.")
		log.Fatal(err)
	}

	usersScanner := bufio.NewScanner(usersFile)

	var users []string

	//goes through users.txt file, store all user info except mine
	for usersScanner.Scan() {
		curUserPN := usersScanner.Text()
		userArr := strings.SplitAfter(curUserPN,",")
		curUser := userArr[0][0:len(userArr[0])-1]
		if user != curUser {
			users = append(users, curUserPN)
		}
	}

	//delete entire users.txt file
	err = os.Remove(usersFileName)
	if err != nil {
		log.Println("Could not delete users file properly.")
	}

	//create users.txt file again
	newUsersFile, err := os.Create(usersFileName)
	if err != nil {
		log.Println("Could not create new users file properly.")
	}

	newUsersFile.Close()

	newUsersFile, err = os.OpenFile(usersFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}

	defer newUsersFile.Close()

	//put back contents of users.txt file, which does not include my info anymore
	putInFile := ""
	for i:=0; i < len(users); i++ {
		putInFile = users[i] + "\n"
		if _, err = newUsersFile.WriteString(putInFile); err != nil {
		    panic(err)
		}
	}
}

func bremoveAcct(user string, lock *Lock) {

	//acquires exclusive lock to write to users.txt, posts.txt and friends.txt
	//once all of these are done, then release
	acquireLock("write", lock)

	removeMeFriendFile(user)
	removeMyPosts(user)
	removeMeUserFile(user)

	releaseLock("write",lock)

}