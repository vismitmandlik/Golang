package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type IpCredential struct {
	Ip          string       `json:"ip"`
	Port        string       `json:"port"`
	Credentials []Credential `json:"credentials"`
}

func main() {
	if len(os.Args) < 2 { // Expecting at least 1 argument: <ipCredentialJson>
		fmt.Println("Usage: go run ssh_command.go <ipCredentialJson>")
		return
	}

	ipCredentialJson := os.Args[1]

	var ipCredentials []IpCredential
	err := json.Unmarshal([]byte(ipCredentialJson), &ipCredentials)
	if err != nil || len(ipCredentials) == 0 { // Check for errors in parsing
		fmt.Println("Invalid or empty IP-Credential data provided")
		return
	}

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

	for _, ipCred := range ipCredentials { // Iterate over each IP-Credential pair
		for _, cred := range ipCred.Credentials { // Try each credential for this IP
			wg.Add(1) // Increment the WaitGroup counter
			go func(ip string, port string, username string, password string) {
				defer wg.Done() // Decrement the counter when the goroutine completes
				if trySSH(ip, port, username, password) {
					fmt.Printf("Successful SSH connection to %s with user %s\n", ip, username)
				}
			}(ipCred.Ip, ipCred.Port, cred.Username, cred.Password) // Pass parameters to the goroutine
		}
	}

	wg.Wait() // Wait for all goroutines to finish
}

func trySSH(ip string, port string, username string, password string) bool {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password), // Use password authentication
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only; use proper host key checking in production
	}

	addr := fmt.Sprintf("%s:%s", ip, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to dial for %s: %s\n", ip, err)
		return false // Return false if connection fails
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s\n", err)
		return false // Return false if session creation fails
	}
	defer session.Close()

	var b []byte
	if b, err = session.CombinedOutput("whoami"); err != nil { // Run a command on the remote server
		fmt.Printf("Failed to run command: %s\n", err)
		return false // Return false if command execution fails
	}

	fmt.Printf("SSH successful to : %s\nCommand Output:\n%s\n", ip, strings.TrimSpace(string(b))) // Print output without extra newlines
	return true                                                                                   // Return true if SSH was successful and command executed
}
