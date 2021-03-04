package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strconv"
        "strings"
)

var count = 0

func handleConnection(c net.Conn) {
        fmt.Print("Connection established \n")
        for {
                netData, err := bufio.NewReader(c).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }

                temp := strings.TrimSpace(string(netData))
                if temp == "STOP" {
                        break
                }

                fmt.Println("Message Received: ",temp)
                counter := strconv.Itoa(count) + "\n"
                c.Write([]byte("Num client connections: "))  //sends message to client
                c.Write([]byte(string(counter)))
        }
        c.Close()
}


//how to use
//in your command line input: go run TCP_Server.go 1234
//      creates TCP server that listens on port number 1234
//      supports multiple concurrent connections
func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide a port number!")
                return
        }

        PORT := ":" + arguments[1]
        l, err := net.Listen("tcp4", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer l.Close()

        for {
                c, err := l.Accept()
                if err != nil {
                        fmt.Println(err)
                        return
                }
                go handleConnection(c)
                count++
        }
}