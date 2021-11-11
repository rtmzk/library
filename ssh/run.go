package ssh

import (
	"os"
)

func (c *Client) Run(cmd string) {
	session, err := c.client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	err = session.Start(cmd)
	if err != nil {
		return
	}

	err = session.Wait()
	if err != nil {
		return
	}
}

func (c *Client) RunWithBindStd(cmd string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	_ = session.Run(cmd)

	return nil
}

func (c *Client) Output(cmd string) ([]byte, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(cmd)
}
