package vault

import (
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func (v *vault) CopyFile(destination string) error {
	// TODO: actually accept a plugin file and io.Copy that instead
	sftpClient, close, err := newSFTPClient(v.Config.SSHAddress, "testuser", "password")
	if err != nil {
		return err
	}
	defer close()

	// Delete file if it exists already, otherwise create a new file
	dstFile, err := sftpClient.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.WriteString(dstFile, "Random text here!\n")
	if err != nil {
		return err
	}

	return nil
}

func newSFTPClient(address, username, password string) (*sftp.Client, func(), error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, nil, err
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	closeConns := func() {
		conn.Close()
		sftpClient.Close()
	}

	return sftpClient, closeConns, nil
}
