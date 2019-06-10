package main

import (
	"bufio"
	"os"
	"strings"
	"log"
)

//checks if a user is already friends with another user
func alreadyFriend(userFriend string , lock *Lock) bool {
	
	//acquires shared lock to read friends.txt
	acquireLock("read", lock)
	
	friendFile, err := os.Open(friendsFileName)
	
	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}
	
	friendScanner := bufio.NewScanner(friendFile)

	defer releaseLock("read",lock)
	defer friendFile.Close()

	//goes through friends file, if there is user,friend pair, return true
	//else return false
	for friendScanner.Scan() {
		curUserFriend := friendScanner.Text()
		if userFriend == curUserFriend {
			return true
		}
	}
	return false
}

//follow someone, write into the friends file user,friend relationship
func writeToFriendFile(userFriend string, lock *Lock) {
	
	userFriend = userFriend + "\n"

	//acquires exclusive lock to write to friends.txt
	acquireLock("write", lock)
	
	friendFile, err := os.OpenFile(friendsFileName, os.O_APPEND|os.O_WRONLY, 0600)
	
	if err != nil {
	    panic(err)
	}

	defer releaseLock("write",lock)
	defer friendFile.Close()

	if _, err = friendFile.WriteString(userFriend); err != nil {
	    panic(err)
	}
}

//takes a user,friend request from FrontEnd, checks if it is allowed, if
//allowed, write to friends file, else return an error message
func addFriend(userFriend string, lock *Lock) (string, bool) {
	
	userFriendArr := strings.SplitAfter(userFriend, ",")
	user := string(userFriendArr[0][0:len(userFriendArr[0])-1])
	friend := string(userFriendArr[1])
	
	//checks if the friend exists
	if checkUserExists(friend, lock) {
		
		//check if the friend is yourself
		if user == friend {
			return "You can't follow yourself!",false
		} else {

			//check if you are already friends
			if alreadyFriend(userFriend, lock) {
				m:= "You are already following " + friend +"!"
				return m,false

			//has to be a new friend, write to friends file
			} else {
				writeToFriendFile(userFriend, lock)
				m:= "You are now following " + friend
				return m,true
			}
		}

	//friend doesn't exist
	} else {
		return "User doesn't exist!",false
	}
}
