package ssh

import (
	"bytes"
	"github.com/pkg/sftp"
	"go-ecm/utils"
	"io"
	"os"
	"path/filepath"
)

func (c *Client) NewSftpClient() *sftp.Client {
	cli, _ := sftp.NewClient(c.client)

	return cli
}

func (c *Client) Upload(src, remote string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return c.UploadDir(src, remote)
	}

	return c.UploadFile(src, remote)
}

func (c *Client) UploadFile(src, remote string) error {
	fileInfo, _ := os.Stat(src)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	var remoteFile, remoteDir string

	if remote[len(remote)-1] == '/' {
		remoteFile = filepath.ToSlash(filepath.Join(remote, filepath.Base(src)))
		remoteDir = remote
	} else {
		remoteFile = remote
		remoteDir = filepath.ToSlash(filepath.Dir(remoteFile))
	}

	if fileInfo.Size() > 1000 {
		remoteSum := c.RemoteMd5Check(remoteFile)
		if remoteSum != "" {
			localSum, _ := utils.Md5File(src)
			if localSum == remoteSum {
				return nil
			}
		}
	}

	client := c.NewSftpClient()
	if err != nil {
		return err
	}
	if client.Stat(remoteDir); err != nil {
		_ = c.RemoteMkdirAll(remoteDir)
	}

	r, _ := client.Create(remoteFile)

	_, err = io.Copy(r, f)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoteMd5Check(path string) string {
	client := c.NewSftpClient()
	_, err := client.Stat(path)
	if err != nil {
		return ""
	}
	b, err := c.Output("md5sum " + path)
	if err != nil {
		return ""
	}
	return string(bytes.Split(b, []byte{' '})[0])
}

func (c *Client) RemoteMkdirAll(path string) error {
	client := c.NewSftpClient()
	err := client.MkdirAll(path)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UploadDir(src, remote string) error { return nil }
