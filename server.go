package main

import (
    "crypto/tls"
    "crypto/rand"
    "time"
    "bufio"
    "fmt"
    "net" 
    "sync" 
    "os"
    "github.com/alexflint/go-arg"
) 

func main() {
    var args struct {
        Host string `arg:"-h" default:"localhost" help:"host/ip to listen on"`
        Port string `arg:"-p" default:"1337" help:"port to listen on"`
        Password string `arg:"required" help:"password required for client connection"`
    }
    arg.MustParse(&args)

    fmt.Println("When you want to distribute commands, just type or paste them here.")

    work := make(chan string)
    go func() {
        s := bufio.NewScanner(os.Stdin)
        for s.Scan() {
            work <- s.Text()
        }
        //close(work) //we want to keep this open to listen for new jobs
    }()

    wg := &sync.WaitGroup{}

    // Set up TLS listener.
    cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}}
    now := time.Now()
    config.Time = func() time.Time { return now }
    config.Rand = rand.Reader
    l, err := tls.Listen("tcp", args.Host + ":" + args.Port, &config)
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on", args.Host, ":", args.Port)

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting:", err.Error())
            os.Exit(1)
        }

        fmt.Println("Received connection:", conn.RemoteAddr().String())
        wg.Add(1)
        r := bufio.NewReader(conn)
        go handleRequest(work, wg, conn, *r, args.Password)
    }
    wg.Wait()
}

func handlePing(conn net.Conn, isDead chan bool, r bufio.Reader) {
    c := make(chan bool)
    go func () {
        for {
            select {
                case ping := <-c:
                    if ping {
                        continue
                    }
                case <-time.After(2 * time.Second):
                    fmt.Println("Closing connection:", conn.RemoteAddr().String())
                    isDead <- true
                    conn.Close()
                    return
            }
        }
    }()
    for {
        time.Sleep( 1 * time.Second )
        message, _ := r.ReadString('\n')
        if message == "PINGPINGPING\n" {
            c <- true     
        } else {
            return
        }
    } 
}

func checkPassword(r bufio.Reader, password string) bool{
    message, _ := r.ReadString('\n')
    if message == password + string('\n') {
        return true
    } else {
        return false
    }
}

func handleRequest(work chan string, wg *sync.WaitGroup, conn net.Conn, r bufio.Reader, password string) {
    defer wg.Done()
    isDead := make(chan bool)
    go handlePing(conn, isDead, r)
    if !checkPassword(r, password){
        fmt.Println("Incorrect password attempt:",conn.RemoteAddr().String())
        conn.Close()
        return
    }
    for text := range work {
        select {
        case itsDead, ok := <-isDead:
            if ok {
                if itsDead {
                    return
                }
            } else {
                fmt.Println("Channel closed!")
                return
            }
        default:
            break
        }
        conn.Write([]byte(text + string('\n')))
        // print a message saying what the message was and where it went
        fmt.Println("Sending to", conn.RemoteAddr().String(), ":", text)
    }
}
