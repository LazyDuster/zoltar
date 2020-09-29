package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var fortunes []string

func fortuneSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '%' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}

func parseFortune(f *os.File) {
	scanner := bufio.NewScanner(f)
	scanner.Split(fortuneSplit)
	for scanner.Scan() {
		fortunes = append(fortunes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading input:%v\n", err)
	}
}

func getFortune() (fortune string) {
	i := rand.Intn(len(fortunes))
	return fortunes[i]
}

func sendFortune(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!fortune" {
		s.ChannelMessageSend(m.ChannelID, getFortune())
	}
}

func main() {
	fmt.Printf("IT'S YOUR LUCKY DAY\n")

	f, err := os.Open("fortunes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading fortune file: %v\n", err)
		os.Exit(1)
	}
	parseFortune(f)
	f.Close()

	if len(os.Args) != 2 {
		fmt.Printf("Usage: gortune [token]\n")
		os.Exit(0)
	}

	var token string = os.Args[1]

	// Create Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Discord session: %v\n", err)
		os.Exit(1)
	}

	dg.AddHandler(sendFortune)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open websocket and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening Discord session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening for commands...\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Printf("\nInterrupt Received! Exiting.\n")
	dg.Close()
}
