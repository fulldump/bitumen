// An example SFTP server implementation using the golang SSH package.
// Serves the whole filesystem visible to the user, and has a hard-coded username and password,
// so not for real use!
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fulldump/goconfig"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Read: https://blog.gopheracademy.com/advent-2015/ssh-server-in-go/
// Read: https://stackoverflow.com/questions/64104586/use-golang-to-get-rsa-key-the-same-way-openssl-genrsa

func main() {
	err := Listen(":3333")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func Listen(address string) error {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			skey := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
			fmt.Println("todo:", skey)
			if skey == "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOAGD6wpG75r1b4R3WFfT0Lgefjt2x/IcpQBNhcN1zHQ" {
				return nil, nil
			}

			return nil, errors.New("UNSUPPORTEDD")
		},
		// PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		// 	// Should use constant-time compare (or better, salt+hash) in
		// 	// a production setting.
		// 	if c.User() == "testuser" && string(pass) == "tiger" {
		// 		return nil, nil
		// 	}
		// 	return nil, fmt.Errorf("password rejected for %q", c.User())
		// },
		BannerCallback: func(conn ssh.ConnMetadata) string {
			fmt.Println("BANNER")
			return "BITUMEN SERVER v0.0.1\n"
		},
		AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
			fmt.Printf("Auth: method '%s', user '%s'\n", method, conn.User())
			if err != nil {
				fmt.Println("ERROR:", err.Error())
			}
		},
		// NoClientAuth: true,
	}

	privateBytes, err := os.ReadFile("hostkey_rsa")
	if err != nil {
		return errors.New("read private key: " + err.Error())
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return errors.New("parse private key: " + err.Error())
	}
	config.AddHostKey(private)

	return listen(config, address)
}

func listen(config *ssh.ServerConfig, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.New("net listen: " + err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			// todo: handle error
			fmt.Println("todo handle err: " + err.Error())
			continue
		}
		sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			// todo: handle error
			fmt.Println("todo handle err: " + err.Error())
			continue
		}
		go ssh.DiscardRequests(reqs)
		go handleServerConn("KeyID????", chans)
		fmt.Println("client version:", string(sConn.ClientVersion()))
	}
}

func handleServerConn(keyID string, chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			fmt.Println("CHANNEL TYPE:", newChan.ChannelType())
			newChan.Reject(ssh.UnknownChannelType, "unkown channel type!")
			continue
		}

		ch, reqs, err := newChan.Accept()
		if err != nil {
			fmt.Println("Handle:" + err.Error())
			continue
		}

		go func(in <-chan *ssh.Request) {
			defer ch.Close()
			for req := range in {

				payload := strings.TrimSpace(string(req.Payload))

				switch req.Type {
				case "exec":
					fmt.Println("EXEC PAYLOAD:", payload)
					ch.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					break
				case "shell":
					fmt.Fprint(ch, "Shell is not allowed\r\n")
					ch.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					return
				case "env":
					// do nothing
				case "pty-req":
				// do nothing
				case "subsystem":
					fmt.Println("SUBSYSTEEEEEMMMMM!!!")
				default:
					fmt.Println("Type", req.Type, payload)

				}
			}
		}(reqs)
	}
}

type Config struct {
	ReadOnly    bool `usage:"read-only server"`
	DebugStdErr bool `usage:"debug to stderr"`
}

// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func main2() {

	c := &Config{
		ReadOnly:    true,
		DebugStdErr: true,
	}
	goconfig.Read(c)

	debugStream := ioutil.Discard
	if c.DebugStdErr {
		debugStream = os.Stderr
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			fmt.Println("password callback")
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Fprintf(debugStream, "Login: %s\n", c.User())
			if c.User() == "testuser" && string(pass) == "tiger" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			fmt.Println("public key callback")
			return nil, nil
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	nConn, err := listener.Accept()
	if err != nil {
		log.Fatal("failed to accept incoming connection", err)
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatal("failed to handshake", err)
	}
	fmt.Fprintf(debugStream, "SSH server established\n")

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatal("could not accept channel.", err)
		}
		fmt.Fprintf(debugStream, "Channel accepted\n")

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.WithDebug(debugStream),
		}

		if c.ReadOnly {
			serverOptions = append(serverOptions, sftp.ReadOnly())
			fmt.Fprintf(debugStream, "Read-only server\n")
		} else {
			fmt.Fprintf(debugStream, "Read write server\n")
		}

		server, err := sftp.NewServer(
			channel,
			serverOptions...,
		)
		if err != nil {
			log.Fatal(err)
		}
		if err := server.Serve(); err == io.EOF {
			server.Close()
			log.Print("sftp client exited session.")
		} else if err != nil {
			log.Fatal("sftp server completed with error:", err)
		}
	}
}
