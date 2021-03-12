package vault

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func (v *vault) WriteFile(sourceFile io.Reader, hostDestination string) error {
	// TODO: parameterisable auth details
	sshConn, err := newSSHClient(v.Config.SSHAddress, "vagrant", "vagrant")
	if err != nil {
		return err
	}

	sftpClient, close, err := newSFTPClient(sshConn)
	if err != nil {
		return err
	}
	defer close()

	// Delete file if it exists already, otherwise create a new file
	dstFile, err := sftpClient.OpenFile(hostDestination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, sourceFile)
	if err != nil {
		return err
	}

	err = makeFileExecutable(dstFile)
	if err != nil {
		return err
	}

	return nil
}

func (v *vault) AddIPCLockCapabilityToFile(filename string) error {
	// TODO: parameterisable auth details
	conn, err := newSSHClient(v.Config.SSHAddress, "vagrant", "vagrant")
	if err != nil {
		return err
	}

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Run(fmt.Sprintf("sudo setcap cap_ipc_lock=ep %s", filename))
	if err != nil {
		return err
	}

	return nil
}

func newSFTPClient(conn *ssh.Client) (*sftp.Client, func(), error) {
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

func newSSHClient(address, username, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func makeFileExecutable(file *sftp.File) error {
	err := file.Chmod(0775)
	if err != nil {
		return err
	}

	return nil
}
