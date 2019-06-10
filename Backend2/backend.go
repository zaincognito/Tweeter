package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
	"strings"
	"sync"
	"time"
	"strconv"
)

var TIMEOUT time.Duration = time.Second
var server_id int
var port string
var usersFileName string
var friendsFileName string
var postsFileName string
var postCountFileName string
var backupServerPort1 string = "localhost:9001"
var backupServerPort2 string = "localhost:9003"


//check if a certain user exists
func checkUserExists(userN string, lock *Lock) bool {
	
	//acquires shared lock tp read file
	acquireLock("read", lock)

	userFile, err := os.Open(usersFileName)

	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	userScanner := bufio.NewScanner(userFile)

	defer releaseLock("read",lock)
	defer userFile.Close()

	//goes through users.txt file, parses each file line into user, checks
	//if user exists, if yes return true, else false
	for userScanner.Scan() {
		curUser := userScanner.Text()
		userArr := strings.SplitAfter(curUser,",")
		if len(userArr) >= 1 && userArr[0] != "" {
			curUsername := string(userArr[0][0:len(userArr[0])-1])
			if userN == curUsername {
				return true
			}
		}
	}

	return false
}

//checks if a username and password pair is valid
func checkValidUser(userN, pass string, lock *Lock) bool {
	
	//acquires shared lock to read file
	acquireLock("read", lock)

	userFile, err := os.Open(usersFileName)
	
	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}
	
	userScanner := bufio.NewScanner(userFile)

	defer releaseLock("read",lock)
	defer userFile.Close()

	//goes through user.txt file, checks each user and password, if match
	//return true, else false
	for userScanner.Scan() {
		curUser := userScanner.Text()
		userArr := strings.SplitAfter(curUser,",") 
		if len(userArr) >= 1 && userArr[0] != "" {
			curUsername := string(userArr[0][0:len(userArr[0])-1])
			curPassword := string(userArr[1][0:len(userArr[1])-1])
			if (userN == curUsername) && (pass == curPassword) {
				return true
			}
		}
	}
	return false
}

//given a username of an existing person, sends back the Name
func getName(username string, lock *Lock) string {
	
	//acquires shared lock to read file
	acquireLock("read", lock)
	
	userFile, err := os.Open(usersFileName)
	
	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}
	
	userScanner := bufio.NewScanner(userFile)

	defer releaseLock("read",lock)
	defer userFile.Close()

	//goes through users.txt file, when current user found return their name
	for userScanner.Scan() {
		curUser := userScanner.Text()
		userArr := strings.SplitAfter(curUser,",")
		curUsername := string(userArr[0][0:len(userArr[0])-1])
		if (username == curUsername) {
			curName := string(userArr[2])
			return curName
		}
	}
	return ""
}

//takes an array of user name, converts to one string with | as separator
func friendsToOneStr(friends []string) string {
	
	if len(friends) == 0 {return "No friends"}
	
	friendStr := friends[0]
	
	for i:= 1; i < len(friends); i++ {
		friendStr = friendStr + " | "
		friendStr = friendStr + friends[i]
	}
	return friendStr
}

func getTransactionNum(port string) int {
	conn, err := net.DialTimeout("tcp", port, TIMEOUT)
	
	if err != nil {
		log.Println("RETURNINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG 0")
		return 0
	} else {
		//putting comma here to ensure the splitting later is working
		s := "transactionNum,"
		fmt.Fprintf(conn, "%s\n",s)

		transactionNum, err1 := bufio.NewReader(conn).ReadString('\n')

		if err1 != nil {
			log.Println("RETURNINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG 0")
			return 0
		} else {
			log.Println("RETURNINGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG 1 tNum")
			log.Println(transactionNum)
			transactionNum = string(transactionNum)
			tNum := 0
			if len(transactionNum) >= 1 && string(transactionNum[0]) != "" {
				transactionNum = transactionNum[0:len(transactionNum)-1]
				tNum,_ = strconv.Atoi(transactionNum)
				log.Println("Helloooooooooooooooooooooooooooooooooooo 1 tNumInt")
				log.Println(tNum)
			}
			return tNum
		}
	}
}

func fetchAllFiles(lock *Lock) string {
	acquireLock("read",lock)

	allFileInfo := ""

	//users.txt
	allFileInfo += usersFileName + "|"
	userFile, err := os.Open(usersFileName)

	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	userScanner := bufio.NewScanner(userFile)
	for userScanner.Scan() {
		line := userScanner.Text()
		allFileInfo += line + ";"

	}
	
	if string(allFileInfo[len(allFileInfo)-1]) != "|" {
		allFileInfo = allFileInfo[0:len(allFileInfo)-1]
	}
	allFileInfo += "~"


	//friends.txt
	allFileInfo += friendsFileName + "|"
	friendsFile, err := os.Open(friendsFileName)

	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	friendScanner := bufio.NewScanner(friendsFile)
	for friendScanner.Scan() {
		line := friendScanner.Text()
		allFileInfo += line + ";"

	}
	//denotes end of file
	if string(allFileInfo[len(allFileInfo)-1]) != "|" {
		allFileInfo = allFileInfo[0:len(allFileInfo)-1]
	}
	allFileInfo += "~"


	//posts.txt
	allFileInfo += postsFileName + "|"
	postsFile, err := os.Open(postsFileName)

	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	postScanner := bufio.NewScanner(postsFile)
	for postScanner.Scan() {
		line := postScanner.Text()
		allFileInfo += line + ";"

	}
	//denotes end of file
	if string(allFileInfo[len(allFileInfo)-1]) != "|" {
		allFileInfo = allFileInfo[0:len(allFileInfo)-1]
	}
	allFileInfo += "~"


	//postCount.txt
	allFileInfo += postCountFileName + "|"
	postCountFile, err := os.Open(postCountFileName)

	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	postCountScanner := bufio.NewScanner(postCountFile)
	for postCountScanner.Scan() {
		line := postCountScanner.Text()
		allFileInfo += line + ";"

	}
	
	if string(allFileInfo[len(allFileInfo)-1]) != "|" {
		allFileInfo = allFileInfo[0:len(allFileInfo)-1]
	}
	allFileInfo += "~"


	releaseLock("read",lock)

	return allFileInfo
}

func updateServer(updatedServerPort string) {
	conn, err := net.Dial("tcp",updatedServerPort)

	if err != nil {
		log.Println("updatedServer not responding")
	}

	s := "giveMeAllUpdates,"

	fmt.Fprintf(conn, "%s\n",s)

	allFileInfo, err1 := bufio.NewReader(conn).ReadString('\n')
			
	if err1 != nil {
		log.Println("Couldnt get all info from up to date server")
	}


	databaseFiles := strings.SplitAfter(allFileInfo, "~")
	databaseFiles = databaseFiles[0:len(databaseFiles)-1]
    for i:=0; i < len(databaseFiles); i++{
        //remove last char for formatting
        databaseFiles[i] = databaseFiles[i][0:len(databaseFiles[i])-1]


        fileNameSplit := strings.SplitAfter(databaseFiles[i], "|")
        //remove last char for formatting
        filePath := fileNameSplit[0][0:len(fileNameSplit[0])-1]
        filePathSplit := strings.SplitAfter(filePath, "/")
	    fileName := filePathSplit[len(filePathSplit)-1]
        
        if string(fileName[0]) == "u"{
        	fileName = usersFileName
        } else if string(fileName[0]) == "f"{
        	fileName = friendsFileName
        } else{
        	if string(fileName[0:5]) == "posts"{
        		fileName = postsFileName
        	} else {
        		fileName = postCountFileName
        	}

        }

        entireFileContent := strings.SplitAfter(fileNameSplit[1], ";")

        for j:=0; j < len(entireFileContent)-1; j++{
            //remove last char for formatting
            entireFileContent[j] = entireFileContent[j][0: len(entireFileContent[j])-1]
        }

        fmt.Println(entireFileContent)            

        err := os.Remove(fileName)
        if err != nil{
            log.Println("couldnt remove %s", fileName)
        }

        newFile, err := os.Create(fileName)
        newFile.Close()

        newFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil{
            panic(err)
        }

        putInFile := ""

        for k:=0; k < len(entireFileContent); k++{
            putInFile = entireFileContent[k] + "\n"
            if _, err = newFile.WriteString(putInFile); err != nil {
                panic(err)
            }
        }

        newFile.Close()

    }

}

func makeMaster(lock *Lock, transaction_num *int) string {
	acquireLock("write",lock)

	numTrans2 := getTransactionNum(backupServerPort1)
	log.Println("numTrans 2 is: ",numTrans2)

	numTrans3 := getTransactionNum(backupServerPort2)
	log.Println("numTrans 3 is: ",numTrans3)

	defer releaseLock("write",lock)
	if *transaction_num >= numTrans3 && *transaction_num >= numTrans2 {
		return "1"
	} else {
		if numTrans2 >= numTrans3 {
			updateServer(backupServerPort1)
			return "1"
		} else {
			updateServer(backupServerPort2)
			return "1"
		}
	}
}

func updateBackUp(backupServerPort, message string, lock *Lock, transaction_num *int) {
	backupConn, err := net.DialTimeout("tcp",backupServerPort,TIMEOUT)

	if err!= nil {
		log.Println("BACKUP SERVER IN PORT %s IS DOWN",backupServerPort)
		return
	}


	pp := getTransactionNum(backupServerPort)
	log.Println("HELLOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO PP: ",pp)
	if pp + 1 < *transaction_num {
		backupConn, err = net.DialTimeout("tcp",backupServerPort,TIMEOUT)

		if err!= nil {
			return
		}
		s := "updateBackUp," + fetchAllFiles(lock)
		log.Println("--------------------")
		log.Println(s)
		log.Println("--------------------")
		fmt.Fprintf(backupConn,"%s\n",s)

	} else {
		backupConn, err = net.DialTimeout("tcp",backupServerPort,TIMEOUT)

		if err!= nil {
			return
		}
		log.Println(message)
		fmt.Fprintf(backupConn,"%s\n",message)
	}

}

//takes a frontend's request to the server in message form, parses the command,
//redirects into the desired function, sends back reply
func handleConn(conn net.Conn, lock *Lock, transaction_num *int) {
	
	scanner:= bufio.NewScanner(conn)

	for scanner.Scan() {
		
		message := scanner.Text()
		fmt.Println(message)
		fmt.Println("-----")

		infoSlice := strings.SplitAfterN(message, ",", 2)
		redirectFunc := string(infoSlice[0][0:len(infoSlice[0])-1])

		if redirectFunc == "signin" {

			userInfo := infoSlice[1]
			exists := bsignin(userInfo, lock)
			var confirmation string
			if exists {
				confirmation = "Yes"
			} else {
				confirmation = "No"
			}
			conn.Write([]byte(confirmation + "\n"))

		} else if redirectFunc == "signup" {
			newUserInfo := infoSlice[1]
			var confirmation string
			created := bsignup(newUserInfo, lock)
			if created {
				confirmation = "Yes"
			} else {
				confirmation = "No"
			}
			updateBackUp(backupServerPort1,"backupsignup,"+newUserInfo,lock,transaction_num)
			updateBackUp(backupServerPort2,"backupsignup,"+newUserInfo,lock,transaction_num)
			*transaction_num++
			conn.Write([]byte(confirmation + "\n"))

		} else if redirectFunc == "usersName" {

			usersName := infoSlice[1]
			name := getName(usersName, lock)
			conn.Write([]byte(name + "\n"))

		} else if redirectFunc == "addFriend" {

			userFriend := infoSlice[1]
			addFriendMessage, _ := addFriend(userFriend, lock)
			updateBackUp(backupServerPort1,"backupaddFriend,"+userFriend,lock,transaction_num)
			updateBackUp(backupServerPort2,"backupaddFriend,"+userFriend,lock,transaction_num)
			*transaction_num++
			conn.Write([]byte(addFriendMessage + "\n"))


		} else if redirectFunc == "writePost" {

			postInfo := infoSlice[1]
			writePost(postInfo, lock)
			updateBackUp(backupServerPort1,"backupwritePost,"+postInfo,lock,transaction_num)
			updateBackUp(backupServerPort2,"backupwritePost,"+postInfo,lock,transaction_num)
			*transaction_num++
			conn.Write([]byte("Yes\n"))

		} else if redirectFunc == "retrievePosts" {

			parameters := strings.SplitAfter(infoSlice[1], ",")
			log.Println("RETREIVEPOSTS ON SERVER %i: %s",server_id,parameters)
			page := parameters[0][0:len(parameters[0])-1]
			username := parameters[1]
			allPosts := retrievePosts(page, username, lock)
			conn.Write([]byte(allPosts + "\n"))

		} else if redirectFunc == "getFriends" {

			user := infoSlice[1]
			friends := getFriends(user, lock)
			friendString := friendsToOneStr(friends)
			conn.Write([]byte(friendString + "\n"))

		} else if redirectFunc == "remFriend" {

			userFriend := infoSlice[1]
			remFriendMessage := removeFriend(userFriend, lock)
			updateBackUp(backupServerPort1,"backupremFriend,"+userFriend,lock,transaction_num)
			updateBackUp(backupServerPort2,"backupremFriend,"+userFriend,lock,transaction_num)
			*transaction_num++
			conn.Write([]byte(remFriendMessage + "\n"))

		} else if redirectFunc == "removeAccount" {

			deleteUser := infoSlice[1]
			bremoveAcct(deleteUser, lock)
			updateBackUp(backupServerPort1,"backupremoveAccount,"+deleteUser,lock,transaction_num)
			updateBackUp(backupServerPort2,"backupremoveAccount,"+deleteUser,lock,transaction_num)
			*transaction_num++
			conn.Write([]byte("Yes\n"))		

		} else if redirectFunc == "checkMaster" {

			serverID := makeMaster(lock,transaction_num)
			conn.Write([]byte(serverID + "\n"))

		} else if redirectFunc == "transactionNum" {

			strTransNum := strconv.Itoa(*transaction_num)
			conn.Write([]byte(strTransNum + "\n"))

		} else if redirectFunc == "giveMeAllUpdates" {

			allFileInfo := fetchAllFiles(lock)
			conn.Write([]byte(allFileInfo + "\n"))

		} else if redirectFunc == "updateBackUp" {
			log.Println("infoSlice[1]: ", infoSlice[1])

			databaseFiles := strings.SplitAfter(infoSlice[1], "~")
			log.Println("DATABASE FILES: ", databaseFiles)
			databaseFiles = databaseFiles[0:len(databaseFiles)-1]
			log.Println("DATABASE FILES AFTER REMOVAL: ", databaseFiles)
		    for i:=0; i < len(databaseFiles); i++{
	            //remove last char for formatting
	            databaseFiles[i] = databaseFiles[i][0:len(databaseFiles[i])-1]

	            fileNameSplit := strings.SplitAfter(databaseFiles[i], "|")
	            //remove last char for formatting
	            filePath := fileNameSplit[0][0:len(fileNameSplit[0])-1]
						            
	            filePathSplit := strings.SplitAfter(filePath, "/")
	            fileName := filePathSplit[len(filePathSplit)-1]
	            // fmt.Println(fileName)
	            if string(fileName[0]) == "u"{
	            	fileName = usersFileName
	            } else if string(fileName[0]) == "f"{
	            	fileName = friendsFileName
	            } else{
	            	if string(fileName[0:5]) == "posts"{
	            		fileName = postsFileName
	            	} else {
	            		fileName = postCountFileName
	            	}

	            }


	            entireFileContent := strings.SplitAfter(fileNameSplit[1], ";")

	            for j:=0; j < len(entireFileContent)-1; j++{
	                //remove last char for formatting
	                entireFileContent[j] = entireFileContent[j][0: len(entireFileContent[j])-1]
	            }


	            err := os.Remove(fileName)
	            if err != nil{
	                log.Println("couldnt remove %s", fileName)
	            }

	            newFile, err := os.Create(fileName)
	            newFile.Close()

	            newFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	            if err != nil{
	                panic(err)
	            }

	            putInFile := ""

	            for k:=0; k < len(entireFileContent); k++{
	                putInFile = entireFileContent[k] + "\n"
	                if _, err = newFile.WriteString(putInFile); err != nil {
	                    log.Println("This is where its breaking")
	                    panic(err)
	                }
	            }
	            log.Println("putInfile: ",putInFile)
	            newFile.Close()
		    }

		} else if redirectFunc == "backupsignup" {
			
			log.Println("transaction_num updated for server: %d -> trans num: %d",server_id,transaction_num)
			newUserInfo := infoSlice[1]
			bsignup(newUserInfo, lock)
			*transaction_num++

		} else if redirectFunc == "backupaddFriend" {

			userFriend := infoSlice[1]
			_, _ = addFriend(userFriend, lock)
			*transaction_num++

		} else if redirectFunc == "backupwritePost" {

			postInfo := infoSlice[1]
			writePost(postInfo, lock)
			*transaction_num++

		} else if redirectFunc == "backupremFriend" {

			userFriend := infoSlice[1]
			removeFriend(userFriend, lock)
			*transaction_num++

		} else if redirectFunc == "backupremoveAccount" {

			deleteUser := infoSlice[1]
			bremoveAcct(deleteUser, lock)
			*transaction_num++

		}

		fmt.Fprintln(os.Stderr, "connection gone!")
		conn.Close()
	}
}

func main() {


	//set up server configurations

	configFile, err := os.Open("server_config.txt")
	if err != nil {
		log.Println("Could not open file properly.")
		log.Fatal(err)
	}

	var configFileLines []string

	configScanner := bufio.NewScanner(configFile)

	for configScanner.Scan() {
		line := configScanner.Text()
		configFileLines = append(configFileLines, line)
	}

	//pull server_id	
	serverIDLine := strings.SplitAfter(configFileLines[0], " ")
	server_id, _ = strconv.Atoi(serverIDLine[len(serverIDLine)-1])

	//pull port number
	portLine := strings.SplitAfter(configFileLines[1], " ")
	port = ":" + portLine[len(portLine)-1]


	//pull database files

	fileLine := strings.SplitAfter(configFileLines[2], " ")
	usersFileName = "../Database/" + fileLine[1][0:len(fileLine[1])-2]
	friendsFileName = "../Database/" + fileLine[2][0:len(fileLine[2])-2]
	postsFileName = "../Database/" + fileLine[3][0:len(fileLine[3])-2]
	postCountFileName = "../Database/" + fileLine[4] //last element doesnt have comma

	//a read/write lock which will be used to lock critical regions in file system
	lock := Lock{sync.RWMutex{}}
	
	transaction_num := 0

	log.Println(port)

	ln, err := net.Listen("tcp", port)
	
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to listen on server")
	}

	t1 := getTransactionNum(backupServerPort1)
	t2 := getTransactionNum(backupServerPort2)
	if t1 != 0 || t2!= 0 {
		if t1 >= t2 {
			updateServer(backupServerPort1)
		} else {
			updateServer(backupServerPort2)
		}
	}
	defer ln.Close()

	for {
		//accepts client's request to connect
		conn, err := ln.Accept()
		
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed to accept")
			os.Exit(1)
		
		}
		fmt.Fprintln(os.Stderr, "accept successful")
		
		//creates thread to process a request, can now process other requests
		//a pass by reference for the read/write lock passed, to protect shared
		//resources
		go handleConn(conn, &lock, &transaction_num)
		
	}
}