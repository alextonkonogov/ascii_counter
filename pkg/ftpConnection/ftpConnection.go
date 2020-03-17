package ftpConnection

import (
	"path/filepath"
	"regexp"

	"github.com/jlaffaye/ftp"
)

// FtpConnection ...
type FtpConnection interface {
	Connect() (err error)
	Disconnect() (err error)
	GetRemoteTxtFileNames(dir string) (txtFiles []string, err error)
	GetConnection() (connection *ftp.ServerConn)
}

type ftpConnection struct {
	host,
	port,
	user,
	password string
	connection *ftp.ServerConn
}

func (f *ftpConnection) Connect() (err error) {
	f.connection, err = ftp.Dial(f.host + ":" + f.port)
	if err != nil {
		return
	}
	err = f.connection.Login(f.user, f.password)

	return
}

func (f *ftpConnection) Disconnect() (err error) {
	err = f.connection.Quit()
	if err != nil {
		return
	}
	err = f.connection.Logout()
	return
}

func (f *ftpConnection) GetConnection() (connection *ftp.ServerConn) {
	return f.connection
}

func (f *ftpConnection) GetRemoteTxtFileNames(dir string) (txtFiles []string, err error) {
	root, err := f.connection.CurrentDir()
	if err != nil {
		return
	}

	list, err := f.connection.List(filepath.Join(root, dir))
	if err != nil {
		return
	}

	for _, v := range list {
		if txt := regexp.MustCompile(`(?m).txt`).MatchString(v.Name); txt {
			txtFiles = append(txtFiles, v.Name)
		}
	}
	return
}

// NewFtp ...
func NewFtp(host, port, user, password string) FtpConnection {
	return &ftpConnection{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}
