package main

import (
	"strings"
)

//checks if a user exists
func bsignin(userInf string, lock *Lock) bool {
	
	potUser := strings.SplitAfterN(userInf, ",", 2)
	userN := string(potUser[0][0:len(potUser[0])-1])
	pass := string(potUser[1])
	isValid := checkValidUser(userN,pass,lock)
	return isValid
}