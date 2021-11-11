package ssh

import (
	"fmt"
	"github.com/zhuihua/library/utils"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

const DefaultSSHPort = 22
const DefaultTimeout = 30 * time.Second

type Client struct {
	Host             string
	Port             int
	AuthMode         int
	Username         string
	Password         string
	KeyFile          string
	HandshakeTimeout time.Duration
	Timeout          time.Duration
	client           *ssh.Client
}

type Option func(c *Client)

func NewClient(username string, port int, options ...Option) (*Client, error) {
	c := &Client{
		Username: username,
		Port: port,
	}

	for _, o := range options {
		o(c)
	}

	if err := c.buildClient(); err != nil {
		return nil, err
	}

	return c, nil
}

func WithAuthByPass(password string) Option {
	return func(c *Client) {
		c.Password = password
	}
}

func WithAuthByKey(key string) Option {
	return func(c *Client) {
		c.AuthMode = 1
		c.KeyFile = key
	}
}

func WithHost(host string) Option {
	return func(c *Client) {
		c.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Client) {
		c.Port = port
	}
}

func WithPassword(pass string) Option {
	return func(c *Client) {
		c.Password = pass
	}
}

func WithTimeout(time time.Duration) Option {
	return func(c *Client) {
		c.Timeout = time
	}
}

func (c *Client) buildClient() error {
	var authFunc ssh.AuthMethod
	switch c.AuthMode {
	case 0:
		authFunc = func() ssh.AuthMethod {
			return passwdAuth(c.Password)
		}()
	case 1:
		authFunc = func() ssh.AuthMethod {
			return keyAuth(c.KeyFile)
		}()
	default:
		return fmt.Errorf("Unknow sshkey auth method")
	}

	cliconf := &ssh.ClientConfig{
		User:            c.Username,
		Timeout:         DefaultTimeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			authFunc,
		},
	}

	if c.Port == 0 {
		c.Port = DefaultSSHPort
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)), cliconf)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func passwdAuth(pw string) ssh.AuthMethod {
	return ssh.Password(pw)
}

func keyAuth(kf string) ssh.AuthMethod {
	if kf == "" {
		kf = utils.UserHome() + "/.ssh/id_rsa"
	}
	key, err := ioutil.ReadFile(kf)
	if err != nil {
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(signer)
}
