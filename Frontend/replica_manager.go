package main

import(
	"net"
	"log"
	"bufio"
	"fmt"
	"strconv"
)

type BEServer struct{
	port string
	server_id int
}

func populateServers(repManagers ServerQueue) ServerQueue {

	//do later: get from file

	server1 := BEServer{port: "localhost:9001",server_id: 1}

	server2 := BEServer{port: "localhost:9002",server_id: 2}

	server3 := BEServer{port: "localhost:9003",server_id: 3}

	repManagers.Push(server1)
	repManagers.Push(server2)
	repManagers.Push(server3)

	return repManagers
}

func chooseNewMaster() int {
	repManagers.Push(repManagers.Front())
	repManagers.Pop()

	newMaster := repManagers.Front()

	service = newMaster.port

	//asks to connect to server
	conn, err := net.DialTimeout("tcp", service,TIMEOUT)
	
	if err != nil {
		log.Println("Server %i is down", newMaster.server_id)
		return chooseNewMaster()
	} else {
		s := "checkMaster,"
		fmt.Fprintf(conn, "%s\n",s)
		
		serverID, err1 := bufio.NewReader(conn).ReadString('\n')
			
		if err1 != nil {
			log.Println("Server %i returned error",newMaster.server_id)
			return chooseNewMaster()
		} else {
			log.Println("Master is server %s",serverID)
			master = newMaster
		}
		ret,_ := strconv.Atoi(serverID)
		return ret
	}
}