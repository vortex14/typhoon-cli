package ssh

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"os"
)

type SSH struct {

}

func (s SSH) TestConnection()  {
	var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: "1",
		Auth: []ssh.AuthMethod{
			ssh.Password("2"),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	client, err := ssh.Dial("tcp", "45.135.233.171:22", config)
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		color.Red("Failed to create session: ", err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		color.Red("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())



}
