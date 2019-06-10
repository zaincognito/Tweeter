package main

import (
	"os"
	"strings"
)

//checks if all signup information is valid, if yes, creates user, else
//returns error message
func bsignup(newUserInf string, lock *Lock) bool {
	
	potNewUser := strings.SplitAfter(newUserInf, ",")
	userN := string(potNewUser[0][0:len(potNewUser[0])-1])
	isValid := checkUserExists(userN, lock)

	//if this username is already used return false
	if isValid {
		return false
	} else {
		
		newUserInf = newUserInf + "\n"

		//acquires exclusive lock to write to users.txt file
		acquireLock("write", lock)
		userFile, err := os.OpenFile(usersFileName, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
		    panic(err)
		}

		defer releaseLock("write",lock)
		defer userFile.Close()

		if _, err = userFile.WriteString(newUserInf); err != nil {
		    panic(err)
		}
		return true
	}
}