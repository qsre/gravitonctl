package cmd

import (
	"fmt"
	"gravitonctl/pkg/aws"
	"gravitonctl/pkg/config"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(connectCmd)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connects to a graviton instance",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Read()
		if err != nil {
			log.Exit(0)
		}

		if len(args) == 0 {
			log.Error("Please supply a name")
			return
		}

		// this code is a mess and needs to be cleaned up!

		key, err := ioutil.ReadFile(c.KeyLocation)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}

		config := &ssh.ClientConfig{
			User: "ec2-user",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			Timeout:         30 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		var ip string

		log.Infof("connecting to instance: %s", args[0])
		for {
			ip, err = aws.GetIp(args[0])
			if err != nil {
				if err.Error() == "public IP isn't available yet" {
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					log.Error(err)
					return
				}
			}

			break
		}

		log.Info(ip)
		hostport := fmt.Sprintf("%s:%d", ip, 22)

		var conn *ssh.Client

		var retries int
		for retries < 5 {

			conn, err = ssh.Dial("tcp", hostport, config)
			if err != nil {
				retries += 1
			} else {
				break
			}

			time.Sleep(1 * time.Second)

		}

		if retries >= 5 {
			log.Errorf("cannot connect to %v: %v", hostport, err)
			return
		}

		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			log.Errorf("cannot open new session: %v", err)
		}
		defer session.Close()

		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			log.Errorf("terminal make raw: %s", err)
		}
		defer terminal.Restore(fd, state)

		w, h, err := terminal.GetSize(fd)
		if err != nil {
			log.Errorf("terminal get size: %s", err)
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		term := os.Getenv("TERM")
		if term == "" {
			term = "xterm-256color"
		}

		if err := session.RequestPty(term, h, w, modes); err != nil {
			log.Errorf("session xterm: %s", err)
		}

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		if err := session.Shell(); err != nil {
			log.Errorf("session shell: %s", err)
		}

		if err := session.Wait(); err != nil {
			if e, ok := err.(*ssh.ExitError); ok {
				switch e.ExitStatus() {
				case 130:
					return
				}
			}
			log.Errorf("ssh: %s", err)
		}
	},
}
