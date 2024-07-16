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

func connectSFTPServer() (*ssh.Client, *sftp.Client) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	SFTPUSER = os.Getenv("SFTP_USER")
	SFTPPSWD = os.Getenv("SFTP_PASSWORD")
	SFTPIP = os.Getenv("SFTP_IP")
	SFTPPATH = os.Getenv("SFTP_PATH")

	config := &ssh.ClientConfig{
		User: SFTPUSER,
		Auth: []ssh.AuthMethod{
			ssh.Password(SFTPPSWD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: fix this insecure
	}

	var client *sftp.Client
	conn, err := ssh.Dial("tcp", SFTPIP, config)
	if err != nil {
		log.Printf("%v\n", NewSFTPError(0, sty.warning.Render("Connection failure. Failed to connect to ssh server.")))
	} else {
		client, err = sftp.NewClient(conn)
		if err != nil {
			log.Printf("%v\n", NewSFTPError(0, sty.warning.Render("Client failure. Failed to create sftp client.")))
			return nil, nil
		}
	}

	return conn, client
}

func closeServer() {
	err := SFTPCLIENT.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%v\n", NewSFTPError(-2, "Failed to close sftp client"))
	}
	err = SSHCLIENT.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%v\n", NewSFTPError(-1, "Failed to close ssh client"))
	}

	fmt.Println(sty.success.Render("Successfully closed connection to server"))
}
