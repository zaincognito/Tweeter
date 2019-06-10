package main

import (
	"bufio"
	"os"
	"strings"
	"log"
)

//remove the relationship where you follow a friend
func removeFriendFile(userFriend string, lock *Lock) {

	//acquires shared lock on friends.txt
	acquireLock("write", lock)

	friendsFile, err := os.Open(friendsFileName)
	if err != nil{
		log.Println("Could not open friends file properly.")
		log.Fatal(err)
	}

	friendsScanner := bufio.NewScanner(friendsFile)

	var friends []string

	//goes through friends.txt file, if line does not have user,friend combo,
	//store into array
	for friendsScanner.Scan() {
		curUserFriend := friendsScanner.Text()
		if userFriend != curUserFriend {
			friends = append(friends, curUserFriend)	
		}
	}

	//entirely delete friends.txt
	err = os.Remove(friendsFileName)
	if err != nil {
		log.Println("Could not delete friends file properly.")
	}

	//creates friends.txt again
	newFriendsFile, err := os.Create(friendsFileName)
	if err != nil {
		log.Println("Could not create new friends file properly.")
	}

	newFriendsFile.Close()

	newFriendsFile, err = os.OpenFile(friendsFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}

	defer releaseLock("write",lock)
	defer newFriendsFile.Close()

	//puts back all friend relationship except user,friend 
	putInFile := ""
	for i:=0; i < len(friends); i++ {
		putInFile = friends[i] + "\n"
		if _, err = newFriendsFile.WriteString(putInFile); err != nil {
		    panic(err)
		}
	}

}

//removes friend from your friends
func removeFriend(userFriend string, lock *Lock) string {
	
	userFriendArr := strings.SplitAfter(userFriend, ",")
	if len(userFriendArr) >= 1 && string(userFriendArr[0]) != "" { 
		user := string(userFriendArr[0][0:len(userFriendArr[0])-1])
		friend := string(userFriendArr[1])
		
		//checks if the friend exists
		if checkUserExists(friend, lock) {

			//checks if friend is yourself
			if user == friend {
				return "You can't unfollow yourself!"
			} else {

				//check if this relation doesn't exists
				if !alreadyFriend(userFriend, lock) {
					m:= "You are not following " + friend +"!"
					return m
				} else {

					//remove the user,friend relationship
					removeFriendFile(userFriend, lock)
					m:= "You have unfollowed " + friend
					return m
				}
			}
		} else {

			//if the user doesn't exist
			return "User doesn't exist!"
		}
	} else { 
		return "User doesn't exist!"
	}
}