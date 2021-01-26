package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"github.com/bwmarrin/discordgo"
)

const (
	words = 5
	commands = 5
)

var fortunes []string
var offensive []string
var commands = []string{"!fortune", "!offendme", "!owoify", "!cloaker", "!fuckme"}
var funcs = []func(){}

/* lol */
func owoify(sentence string) (out string) {
	faces := []string{"owo", "OWO", "OwO", "UwU", ">w<", "^w^", "uwu"}
	glyphs := [words][2]string{
		{"l", "w"},
		{"r", "w"},
		{"L", "W"},
		{"R", "W"},
		{"!", faces[rand.Intn(6)]},
	}
	out = strings.ReplaceAll(sentence, glyphs[0][0], glyphs[0][1])
	for i := 1; i < words; i++ {
		out = strings.ReplaceAll(out, glyphs[i][0], glyphs[i][1])
	}
	return out
}

/* Scanner split function, reads in file and delimits by % */
func FortuneSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

/* Runs through the given file using the splitter function and fills the fortune slices */
func ParseFortune(f *os.File, o bool) {
	scanner := bufio.NewScanner(f)
	scanner.Split(FortuneSplit)
	for scanner.Scan() {
		if o {
			offensive = append(fortunes, scanner.Text())
		} else {
			fortunes = append(fortunes, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading input:%v\n", err)
	}
}

/* Return a fortune randomly */
func GetFortune() (fortune string) {
	i := rand.Intn(len(fortunes))
	return fortunes[i]
}

/* Upon receiving a command, bot sends a fortune */
func SendFortune(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmdstr := strings.SplitN(m.Content, " ", 2)

	if cmdstr[0] == "!fortune" {
		msgtitle := m.Author.Username + ", it's your lucky day."
		me := discordgo.MessageEmbed{ Title: msgtitle, Description: GetFortune(), Color: 39423 }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!offendme" {
		msgtitle := m.Author.Username + ", it's your lucky day."
		me := discordgo.MessageEmbed{ Title: msgtitle, Description: GetFortune(), Color: 14103594 }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!owoify" {
		msgtitle := m.Author.Username + ", are you fucking kidding me?"
		me := discordgo.MessageEmbed{ Title: msgtitle, Description: owoify(cmdstr[1]), Color: 0xFFC0CB }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!fuckme" {
		me := discordgo.MessageEmbed{ Title: "An odd request...", Description: "...but okay, \\*cums\\*.", Color: 0xFFC0CB }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!fuckmeinitalian" {
		me := discordgo.MessageEmbed{ Title: "A spicy meatball...", Description: "...but okay, *cums*.", Color: 0xFFC0CB }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!cloaker" {
		me := discordgo.MessageEmbed{ Title: "WHOOOOOO!", Description: "CALL ME THE CLOAKER SMOKER!", Color: 0x000000 }
		s.ChannelMessageSendEmbed(m.ChannelID, &me)
	} else if cmdstr[0] == "!zoltar" {
		var greeting string = "ZOLTAR SAYS: Make your wish, !fortune."
		s.ChannelMessageSend(m.ChannelID, greeting)
	}
	/*
	for i := 0; i < 6; i++ {
		if cmdstr[0] == commands[0] {
			//
		}
	}
	*/
}

func main() {
	fmt.Printf("IT'S YOUR LUCKY DAY\n")

	f, err := os.Open("fortunes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading fortune file: %v\n", err)
		os.Exit(1)
	}
	ParseFortune(f, false)
	f.Close()

	f, err = os.Open("offensive")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading offensive fortune file: %v\n", err)
		os.Exit(1)
	}
	ParseFortune(f, true)
	f.Close()

	if len(os.Args) != 2 {
		fmt.Printf("Usage: zoltar [token]\n")
		os.Exit(0)
	}

	var token string = os.Args[1]
	rand.Seed(time.Now().UTC().UnixNano())

	// Create Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Discord session: %v\n", err)
		os.Exit(1)
	}

	dg.AddHandler(SendFortune)
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
