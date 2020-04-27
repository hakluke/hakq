package main

import (
    "bufio"
    "fmt"
    "net"
    "os/exec"
    "strings"
    "sync"
    "time"
    "github.com/alexflint/go-arg"
    "crypto/tls"
)


func main() {
    var args struct {
        Threads int     `arg:"-t" default:"5" help:"number of threads"`
        Server string   `arg:"required" help:"server to connect to" placeholder:"IP:PORT"`
        Password string `arg:"required" help:"password for server" placeholder:"PASSWORD"`   
        Insecure bool `help:"ignore insecure TLS - not recommended!"`
    }
    arg.MustParse(&args)
    tlsconfig := tls.Config{}
    if args.Insecure {
        tlsconfig = tls.Config{InsecureSkipVerify: true} 
    }    
    c, err := tls.Dial("tcp", args.Server, &tlsconfig)
    for {
        if err == nil {
            break
        } else {
            fmt.Println("Retrying after error:",err)
        }
        time.Sleep(3000 * time.Millisecond)
        c, err = tls.Dial("tcp", args.Server, nil)
    }

    // set up reader
    r := bufio.NewReader(c)

    // send password
    sendLine(c, args.Password)

    // set up keepalive pings
    go ping(c)

    // set up workers
    wg := &sync.WaitGroup{}
    for i := 1; i <= args.Threads; i++ {
        wg.Add(1)
        go worker(wg, c, i, r)
    }
    
    wg.Wait()
}

func ping(c net.Conn) {
    for {
        c.Write([]byte("PINGPINGPING" + string('\n')))
        time.Sleep(1 * time.Second)
    }
}

func sendLine(c net.Conn, line string) {
    c.Write([]byte(line + string('\n')))
}

func worker(wg *sync.WaitGroup, c net.Conn, workerNumber int, r *bufio.Reader) {
    defer wg.Done()
    defer c.Close()
    for {
        //receive messages
        message, _:= r.ReadString('\n')
        if message == "" {
            c.Close() // this is usually an indication that the connection has dropped
            return
        }
        //trim the newline
        message = strings.TrimSuffix(message, "\n")
        fmt.Println("Worker", workerNumber, "executing:", message)
        //execute the command
        out, err := exec.Command("/bin/sh", "-c", message).Output()
        if err != nil {
                fmt.Println(err)
        }
        //print the command output
        fmt.Printf("%s", out)
    }
}
