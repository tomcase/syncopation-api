package sync

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"github.com/tomcase/syncopation-api/models"
	"golang.org/x/crypto/ssh"
)

func Go(ctx models.ServerService) error {
	response, err := ctx.List(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	if len(response) == 0 {
		return errors.New("there are no servers for sync to run against")
	}

	server := response[0]

	config := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%v", server.Host, server.Port)
	sshConn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("failed to dial ssh: %v", err)
	}
	defer sshConn.Close()

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP Client: %v", err)
	}
	defer client.Close()

	source := server.SourcePath
	err = downloadDir(client, source, server)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	return nil
}

func downloadDir(client *sftp.Client, dirPath string, server *models.Server) error {
	fi, err := client.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to list files at destination '%s': %v", dirPath, err)
	}
	for _, file := range fi {
		fp := path.Join(dirPath, file.Name())
		if file.IsDir() {
			err = downloadDir(client, fp, server)
			if err != nil {
				return err
			}
			continue
		}

		err = downloadFile(client, fp, server)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(client *sftp.Client, dirPath string, server *models.Server) error {
	s := strings.Split(dirPath, string(os.PathSeparator))
	dir := path.Join(server.DestinationPath, path.Join(s[1:len(s)-1]...))
	err := os.MkdirAll("/incomplete", os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create directory at '/incomplete': %v", err)
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create directory at %s: %v", dir, err)
	}

	incPath, err := copyToTempLocation(client, dirPath)
	if err != nil {
		return err
	}
	tmpFile, err := os.Open(incPath)
	if err != nil {
		return fmt.Errorf("could not open temp file at %s: %v", dirPath, err)
	}
	defer tmpFile.Close()

	destPath := path.Join(server.DestinationPath, path.Join(s[1:]...))
	df, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not open file at %s: %v", destPath, err)
	}
	defer df.Close()

	bytes, err := io.Copy(df, tmpFile)
	if err != nil {
		return fmt.Errorf("could not copy remote file at %s: %v", dirPath, err)
	}
	log.Printf("Downloaded %s - %d bytes copied\n", dirPath, bytes)

	err = os.Remove(incPath)
	if err != nil {
		return fmt.Errorf("failed to remove temp file at %s: %v", incPath, err)
	}

	return nil
}

func copyToTempLocation(client *sftp.Client, dirPath string) (string, error) {
	sf, err := client.Open(dirPath)
	if err != nil {
		return "", fmt.Errorf("could not open remote file at %s: %v", dirPath, err)
	}
	defer sf.Close()

	incPath := path.Join("/incomplete", uuid.New().String())
	incFile, err := os.Create(incPath)
	if err != nil {
		return "", fmt.Errorf("could not open temp file at %s: %v", incPath, err)
	}
	defer incFile.Close()

	_, err = io.Copy(incFile, sf)
	if err != nil {
		return "", fmt.Errorf("could not copy remote file at %s: %v", dirPath, err)
	}

	return incPath, nil
}
