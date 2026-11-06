package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	// "golang.org/x/crypto/ssh"
)

type sshInfo struct {
	user     string
	addr     string
	keyFile  string
	command  string
	password string
	out      string
}

func extractSshInfo() sshInfo {
	infos := sshInfo{}
	env, err := os.ReadFile("ssh.env")
	if err != nil {
		log.Fatalf("Could not read enviroment file: %v", err)
	}
	lines := strings.Split(string(env), "\n")
	for i := 0; i < len(lines); i++ {
		split := strings.Split(lines[i], ":")
		if split[0] == "user" {
			infos.user = split[1]
		}
		if split[0] == "server" {
			infos.addr = split[1]
		}
		if split[0] == "command" {
			infos.command = split[1]
		}
		if split[0] == "keyfile" {
			infos.keyFile = split[1]
		}
		if split[0] == "password" {
			infos.password = split[1]
		}
		if split[0] == "out" {
			infos.out = split[1]
		}
	}

	return infos
	// return sshInfo{addr: infos.addr, user: infos.user, keyFile: infos.keyFile, command: infos.command, password: infos.password}
}
func executeSshCommand(info sshInfo) {
	key, err := os.ReadFile(info.keyFile)
	if err != nil {
		log.Fatalf("Could not read file: %v, %s", err, info.keyFile)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Could not parse key: %v", err)
	}
	// Hardcoded port
	addr := net.JoinHostPort(info.addr, strconv.Itoa(22))

	config := &ssh.ClientConfig{
		User: info.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// Change to known host later
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	if info.password != "" {
		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		stdin, err := session.StdinPipe()
		if err != nil {
			log.Fatalf("StdinPipe error: %v", err)
		}

		if err := session.Start(info.command); err != nil {
			log.Fatalf("Start command error: %v", err)
		}
		_, err = io.WriteString(stdin, info.password)
		if err != nil {
			log.Fatalf("Write password error: %v", err)
		}
		stdin.Close()

		if err := session.Wait(); err != nil {
			log.Fatalf("Command failed: %v | Stderr: %s", err, stderrBuf.String())
		}
	} else {
		_, err := session.Output(info.command)
		if err != nil {
			log.Fatalf("Failed to run command: %v", err)
		}
	}

}

func downloadViaCat(info sshInfo, localPath string) error {
	key, err := os.ReadFile(info.keyFile)
	if err != nil {
		log.Fatalf("Could not read file: %v, %s", err, info.keyFile)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Could not parse key: %v", err)
	}
	// Hardcoded port
	addr := net.JoinHostPort(info.addr, strconv.Itoa(22))

	config := &ssh.ClientConfig{
		User: info.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// Change to known host later
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// 1. Create the destination file on the local machine
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// 2. Direct session.Stdout straight into the local file stream
	session.Stdout = localFile

	// 3. Capture stderr in a buffer to catch errors (e.g., file not found, permission denied)
	var errBuf bytes.Buffer
	session.Stderr = &errBuf

	// 4. Run 'cat' on the remote host
	// Using single quotes around the path protects against spaces in filenames
	cmd := fmt.Sprintf("cat '%s'", info.out)
	if err := session.Run(cmd); err != nil {
		// Clean up the local empty/incomplete file if remote read failed
		os.Remove(localPath)
		return fmt.Errorf("remote read failed: %w (stderr: %s)", err, errBuf.String())
	}

	return localFile.Sync()
}

func jsonReader(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data map[string]any
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}

	for key := range data {
		fmt.Println(key)
	}
}
