package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "10.101.171.173:38250")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("!!")
	//wbuffer := make([]byte, 20)
	rbuffer := make([]byte, 2048)
	//fmt.Scanf("%s", &wbuffer)
	if len(os.Args) == 3 {
		conn.Write([]byte(os.Args[1] + " " + os.Args[2]))
	} else if len(os.Args) == 4 {
		conn.Write([]byte(os.Args[1] + " " + os.Args[2] + " " + os.Args[3]))
	} else if len(os.Args) == 5 {
		conn.Write([]byte(os.Args[1] + " " + os.Args[2] + " " + os.Args[3] + " " + os.Args[4]))
	}
	n, err := conn.Read(rbuffer)
	m := n
	for n != 0 {
		n1, err := conn.Read(rbuffer[m:])
		m += +n1
		if n1 == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		n = n1
	}
	fmt.Printf(string(rbuffer[0:m]))
	conn.Close()

}
