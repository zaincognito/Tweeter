package main

import (
	"bufio"
	"os"
	"strings"
	"log"
	"strconv"
	"io/ioutil"
)

//takes a post written in FrontEnd and stores it in posts.txt
func writePost(postInfo string, lock *Lock) {
	
	//acquires shared lock to read postCount.txt
	acquireLock("read", lock)

	postCountF, err := os.Open(postCountFileName)
	
	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}
	
	postCountScanner := bufio.NewScanner(postCountF)
	postCountScanner.Scan()
	postCount := postCountScanner.Text()

	//taking the postCount and making it an int
	postCountI, _ := strconv.Atoi(postCount)

	postCountF.Close()
	releaseLock("read",lock)

	//increment post count, convert it back to a string
	postCountI++;
	log.Println(postCountI)
	theID := strconv.Itoa(postCountI)
	log.Println("string: %s",theID)

	writeCount := []byte(theID)

	//acquires exclusive lock to write into postCount.txt
	acquireLock("write", lock)
	ioutil.WriteFile(postCountFileName,writeCount,0644)
	releaseLock("write",lock)

	//setting up post/post info to store into posts.txt
	theID = theID + ","
	postArr := strings.SplitAfter(postInfo, ",")
	user := postArr[0]
	post := postArr[1]
	thePost := user + theID + post
	thePost = thePost + "\n"
	
	//acquires exclusive lock to write into posts.txt
	acquireLock("write", lock)
	
	postFile, err := os.OpenFile(postsFileName, os.O_APPEND|os.O_WRONLY, 0600)
	
	if err != nil {
	    panic(err)
	}

	defer releaseLock("write",lock)
	defer postFile.Close()

	if _, err = postFile.WriteString(thePost); err != nil {
	    panic(err)
	}

}