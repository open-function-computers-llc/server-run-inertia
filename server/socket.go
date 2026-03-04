package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/open-function-computers-llc/server-run-inertia/account"
	"github.com/sirupsen/logrus"
)

func (s *server) setupSocket() {
	socketPath := os.Getenv("APP_SOCKET_PATH")

	// Remove existing socket file
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		s.logger.Error("Error creating Unix socket:", err)
		return
	}

	err = os.Chmod(socketPath, 0666)
	if err != nil {
		s.logger.Error("Error satanizing Unix socket:", err) // 🤘
		return
	}

	// on SE Linux? make the socket happy by default
	exec.Command("/usr/sbin/restorecon", socketPath).Run()

	defer listener.Close()
	s.logger.Info("Listening on", socketPath)

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting connection:", err)
			continue
		}

		wg.Add(1)
		go handleConnection(conn, &wg, s.defaultLogo, s.logger)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup, defaultLogo []byte, l *logrus.Logger) {
	defer wg.Done()
	defer conn.Close()

	scriptsRoot := os.Getenv("SCRIPTS_ROOT")
	var cmd *exec.Cmd

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read error:", err)
		return
	}

	var msg SocketMessage
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		log.Println("invalid JSON:", err)
		conn.Write([]byte("ERROR: invalid JSON\n"))
		return
	}

	log.Printf("Account: %s, Command: %s\n", msg.Account, msg.Command)

	// Example command handling
	switch msg.Command {
	case "LOCK-STATUS":
		hostingAccount, err := account.Find(msg.Account)
		if err != nil {
			conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		}

		lockStatus := "UNLOCKED"
		if hostingAccount.IsLocked {
			lockStatus = "LOCKED"
		}
		conn.Write([]byte("OK: " + lockStatus + "\n"))
	case "LOGO":
		logoBytes := defaultLogo
		customLogoBytes, err := useCustomLogo()
		if err == nil {
			logoBytes = customLogoBytes
		}

		base64Bytes := base64.StdEncoding.EncodeToString(logoBytes)
		conn.Write([]byte("OK: " + base64Bytes + "\n"))
	case "UNLOCK":
		cmd = exec.Command(scriptsRoot+"ofco-unlock-site.sh", msg.Account)
		cmd.Stderr = cmd.Stdout
		cmd.Env = append(cmd.Env, "NOCONFIRM=yes")
		cmd.Env = append(cmd.Env, "TERM=xterm-mono")

		output, err := cmd.Output()
		if err != nil {
			conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		} else {
			conn.Write([]byte("OK: " + string(output) + "\n"))
		}
	case "LOCK":
		l.Info("trying to lock!")
		l.Info("running command " + scriptsRoot + "ofco-lock-site-production.sh " + msg.Account)
		cmd = exec.Command(scriptsRoot+"ofco-lock-site-production.sh", msg.Account)
		cmd.Stderr = cmd.Stdout
		cmd.Env = append(cmd.Env, "NOCONFIRM=yes")
		cmd.Env = append(cmd.Env, "TERM=xterm-mono")

		output, err := cmd.Output()
		if err != nil {
			conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		} else {
			conn.Write([]byte("OK: " + string(output) + "\n"))
		}
	default:
		conn.Write([]byte("ERROR: unknown command\n"))
	}
}

type SocketMessage struct {
	Account string `json:"account"`
	Command string `json:"command"`
}
