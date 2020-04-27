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
        "strconv"
)


func main() {
    var args struct {
        Threads int     `arg:"-t" default:"5" help:"number of threads"`
        Server string   `arg:"required" help:"server to connect to" placeholder:"IP:PORT"`
        Password string `arg:"required" help:"password for server" placeholder:"PASSWORD"`   
    }
    arg.MustParse(&args)
    
    serverSplit := strings.Split(args.Server, ":")
    serverIP := net.ParseIP(serverSplit[0])
    port, _ := strconv.Atoi(serverSplit[1])
    server := net.TCPAddr{IP: serverIP, Port: port} 
    client := net.TCPAddr{IP: net.ParseIP("127.0.0.1")}
    c, err := net.DialTCP("tcp", &client, &server)
    for {
        if err == nil {
            break
        } else {
            fmt.Println("Retrying after error:",err)
        }
        time.Sleep(3000 * time.Millisecond)
        c, err = net.DialTCP("tcp", &client, &server)
    }

    // set up workers
    wg := &sync.WaitGroup{}
    for i := 1; i <= args.Threads; i++ {
        wg.Add(1)
        go worker(wg, c, i)
    }
    
    wg.Wait()
}

func ping(c net.Conn) {
    for {
        c.Write([]byte("PINGPINGPING" + string('\n')))
        time.Sleep(900 * time.Millisecond)
    }
}

func worker(wg *sync.WaitGroup, c net.Conn, workerNumber int) {
    defer wg.Done()
    go ping(c)
    defer c.Close()
    for {
        //receive messages
        message, _ := bufio.NewReader(c).ReadString('\n')
        if message == "" {
            fmt.Println("No tasks, sleeping.")
            time.Sleep(1000 * time.Millisecond)
            continue
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
