package main

import (
    "bufio"
    "flag"
    "fmt"
    "net" 
    "sync" 
    "os"
) 

const (
    CONN_HOST = "localhost"
    CONN_PORT = 1337
    CONN_TYPE = "tcp"
) 


func main() {
        flag.Parse()
        fmt.Println("When you want to distribute commands, just type or paste them here.")

        work := make(chan string)
        go func() {
            s := bufio.NewScanner(os.Stdin)
            for s.Scan() {
                work <- s.Text()
            }
            //close(work)
        }()

        wg := &sync.WaitGroup{}

        // Listen for incoming connections.
        listenIP := net.ParseIP(CONN_HOST)
        host := net.TCPAddr{IP: listenIP, Port: CONN_PORT}
        l, err := net.ListenTCP(CONN_TYPE, &host)
        if err != nil {
            fmt.Println("Error listening:", err.Error())
            os.Exit(1)
        }

        // Close the listener when the application closes.
        defer l.Close()
        fmt.Println("Listening on", CONN_HOST, ":", CONN_PORT)

        for {
            // Listen for an incoming connection.
            conn, err := l.AcceptTCP()
            if err != nil {
                fmt.Println("Error accepting:", err.Error())
                os.Exit(1)
            }

            fmt.Println("Received connection:", conn.RemoteAddr().String())
            wg.Add(1)
            go handleRequest(work, wg, conn)
        }
        wg.Wait()
}

func handleRequest(work chan string, wg *sync.WaitGroup, conn net.Conn) {
    defer wg.Done()
    for text := range work {
        // write the job to the connection
        conn.Write([]byte(text + string('\n')))
        // print a message saying what the message was and where it went
        fmt.Println("Sending to", conn.RemoteAddr().String(), ":", text)
    }
}
