package main

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "log"
    "os"
    "bufio"
    "strings"
    "path/filepath"
    "errors"
 )

func main() {
    hostKey, err := getHostKey("<ssh target host>")
    if err != nil {
            log.Fatal(err)
    }

    config := &ssh.ClientConfig{
        User: "<ssh user>",
        Auth: []ssh.AuthMethod{
            ssh.Password("<ssh ldap password>"),
            },
        HostKeyCallback: ssh.FixedHostKey(hostKey),
    }

    // Connect to the remote server and perform the SSH handshake.
    sshconn, err := ssh.Dial("tcp", "<ssh target host>:22", config)
    if err != nil {
            log.Fatalf("unable to connect: %v", err)
    }

    defer sshconn.Close()
    fmt.Println("connected!!!")
}

func getHostKey(host string) (ssh.PublicKey, error) {
    file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var hostKey ssh.PublicKey
    for scanner.Scan() {
        fields := strings.Split(scanner.Text(), " ")
        if len(fields) != 3 {
            continue
        }
        if strings.Contains(fields[0], host) {
            var err error
            hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
            if err != nil {
                return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
            }
            break
        }
    }

    if hostKey == nil {
        return nil, errors.New(fmt.Sprintf("no hostkey for %s", host))
    }
    return hostKey, nil
}
