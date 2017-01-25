package main

import (
    "os"
	"golang.org/x/crypto/ssh"
    "io"
    "strings"
)

func main () {
    argsWithoutProg := os.Args[1:]
    
    sshConfig := &ssh.ClientConfig{
        User: "root",
        Auth: []ssh.AuthMethod{
            ssh.Password("root"),
        },
    }
    
    conn, err := ssh.Dial("tcp", "localhost:2222", sshConfig)
    if err != nil {
        panic(err)
    }

    defer conn.Close()

    session, serr := conn.NewSession()
    if serr != nil {
        panic(serr)
    }

    modes := ssh.TerminalModes{
        ssh.ECHO:          0,     // disable echoing
        ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
        ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
    }

    if perr := session.RequestPty("xterm", 80, 40, modes); perr != nil {
        session.Close()
        panic(perr)
    }

    defer session.Close()

    sin, _ := session.StdinPipe()
    go io.Copy(sin, os.Stdin)
    
    sou, _ := session.StdoutPipe()
    go io.Copy(os.Stdout, sou)

    ser, _ := session.StderrPipe()
    go io.Copy(os.Stderr, ser)
    
    args := strings.Join(argsWithoutProg, " ")
    cmd := strings.Join([]string{"/project/pythonenv/bin/python", args}, " ")
    
    rerr := session.Run(cmd)
    if err != nil {
        panic(rerr)
    }


}