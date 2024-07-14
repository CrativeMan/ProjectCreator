package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPError struct {
	Code    int
	Message string
}

func (e *SFTPError) Error() string {
	return fmt.Sprintf("SFTPError - Code: %d, Message: %s", e.Code, e.Message)
}

func NewSFTPError(code int, message string) error {
	return &SFTPError{
		Code:    code,
		Message: message,
	}
}

func createSFTPConnectionClient() (*sftp.Client, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Err loading .env file: %v", err)
	}

	sftpServer := os.Getenv("SFTP_IP")
	sftpUser := os.Getenv("SFTP_USER")
	sftpPassword := os.Getenv("SFTP_PASSWORD")

	fmt.Printf("serv: %v, usr: %v, pswd: %v\n", sftpServer, sftpUser, sftpPassword)

	config := &ssh.ClientConfig{
		User: sftpUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: FIX THIS, INSECURE
	}

	conn, err := ssh.Dial("tcp", sftpServer, config)
	if err != nil {
		return nil, NewSFTPError(0, "Failed to connect to sftpserver, get files localy")
	}
	defer conn.Close()
	return initSFTPClient(conn), err
}

func initSFTPClient(conn *ssh.Client) *sftp.Client {
	client, err := sftp.NewClient(conn)
	if err != nil {
		// TODO: error handling
		fmt.Println("Failed to create SFTP client, getting files locally.")
	}
	defer client.Close()

	return client
}

func getAllRemoteFiles() {
	cwd, err := SFTPCLIENT.Getwd()
	if err != nil {
		panic(err) // TODO: error handling
	}
	fmt.Println(cwd)
}
