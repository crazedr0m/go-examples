package main
import (
	"bufio"
    "fmt"
    "net"
	"os"
)
func main() {
	PORT := ":9090" 
    listener, err := net.Listen("tcp", PORT) 
     
    if err != nil {
        fmt.Println(err) 
        return
    } 
    defer listener.Close() 
    fmt.Println("Server is listening...")

    for { 
        conn, err := listener.Accept() 
        if err != nil { 
            fmt.Println(err) 
			 os.Exit(1)
        } 

		go handleRequest(conn)
    } 
}

func handleRequest(conn net.Conn) {
    defer conn.Close()
    // читаем данные от клиента
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        clientMessage := scanner.Text()
        fmt.Printf("Received from client: %s\n", clientMessage)
        // отправляем ответ клиенту
        conn.Write([]byte("Message received.\n"))
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading:", err.Error())
    }
}