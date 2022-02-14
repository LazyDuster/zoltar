package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"cgt.name/pkg/go-mwclient"
	
	"github.com/bwmarrin/discordgo"
)

const (
	words = 11
	command_len = 3
)

var fortunes []string
var offensive []string

/* Gets uptime for the system. */
func GetUptime() (timevals [3]int) {
	system := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(system); err != nil {
		return timevals
	} else {
		system.Uptime /= 60
		timevals[0] = int(system.Uptime % 60)
		system.Uptime /= 60
		timevals[1] = int(system.Uptime % 24)
		timevals[2] = int(system.Uptime / 24)
		return timevals
	}
}

/* lol */
func owoify(sentence string) (out string) {
	faces := []string{"owo", "OWO", "OwO", "UwU", ">w<", "^w^", "uwu"}
	glyphs := [words][2]string{
		{"l", "w"},
		{"r", "w"},
		{"L", "W"},
		{"R", "W"},
		/* user suggestion: replace all phallic references to "bulgey wulgey" */
		{"penis", "widdwe pickwe"},
		{"PENIS", "WIDDWE PICKWE"},
		{"cock", "bulgie wulgie"},
		{"COCK", "BULGIE WULGIE"},
		{"dick", "bulgie wulgie"},
		{"DICK", "BULGIE WULGIE"},
		{"!", faces[rand.Intn(6)]},
	}
	out = strings.ReplaceAll(sentence, glyphs[0][0], glyphs[0][1])
	for i := 1; i < words; i++ {
		out = strings.ReplaceAll(out, glyphs[i][0], glyphs[i][1])
	}
	return out
}

/* Opens a MediaWiki client and returns a random page. */
func ReturnWiki() (url string, err error) {
	w, err := mwclient.New("https://en.wikipedia.org/w/api.php", "zoltarBot")
	if err != nil {
		return "Problem creating wiki client: " + err.Error(), err
	}

	parameters := map[string]string{
		"action":       "query",
		"prop":         "info",
		"inprop":       "url",
		"generator":    "random",
		"grnnamespace": "0",
		"grnlimit":     "1",
	}

	res, err := w.Get(parameters)
	if err != nil {
		return "Problem making api request: " + err.Error(), err
	}

	temp, _ := res.GetObjectArray("query", "pages")
	return temp[0].GetString("fullurl")
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
			offensive = append(offensive, scanner.Text())
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

/* Return an offensive fortune randomly */
func GetOffensive() (fortune string) {
	i := rand.Intn(len(offensive))
	return offensive[i]
}

/* Upon receiving a command, bot parses and executes */
func ParseCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	var msgtitle string
	var clr int
	var desc string
	cmdstr := strings.SplitN(m.Content, " ", 2)

	switch cmdstr[0] {
	case "!uptime":
		msgtitle = m.Author.Username + ", ZOLTAR has lived since time immemorial!"
		timevals := GetUptime()
		desc += "System uptime is: "
		if timevals[2] != 0 {
			desc += strconv.Itoa(timevals[2]) + " days, "
		}
		if timevals[1] != 0 {
			desc += strconv.Itoa(timevals[1]) + " hours, and "
			desc += strconv.Itoa(timevals[0]) + " minutes."
		} else {
			desc += strconv.Itoa(timevals[0]) + " minutes."
		}
		clr = 0x32CD32
	case "!downtime":
		msgtitle = m.Author.Username + ", the mighty ZOLTAR never goes down!"
		desc = "...and even if I did, it'd be for less than mee6 :^)"
		clr = 0x32CD32
	/* Youtube command, currently not working */
	case "!q":
		msgtitle = m.Author.Username + " wills it, so it shall be done."
		item := SearchYT(cmdstr[1])
		if item == nil {
			s.ChannelMessageSend(m.ChannelID, "Search returned no results.")
			return
		}
		desc = item.Title
		clr = 0xC4302B
	case "!fortune":
		msgtitle = m.Author.Username + ", it's your lucky day."
		desc = GetFortune()
		clr = 0x0099FF
	case "!offendme":
		msgtitle = m.Author.Username + ", it's your lucky day."
		desc = GetOffensive()
		clr = 0xD7342A
	case "!owoify":
		msgtitle = m.Author.Username + ", are you fucking kidding me?"
		desc = owoify(cmdstr[1])
		clr = 0xFFC0CB
	case "!fuckme":
		msgtitle = "An odd request..."
		desc = "...but okay, \\*cums\\*."
		clr = 0xFFC0CB
	case "!fuckmeinitalian":
		msgtitle = "A spicy meatball..."
		desc = "...but okay, *cums*."
		clr = 0xFFC0CB
	case "!cloaker":
		msgtitle = "WHOOOOOO!"
		desc = "CALL ME THE CLOAKER SMOKER!"
		clr = 0x000000
	case "!cytube":
		msgtitle = m.Author.Username + ", a marvelous idea!"
		desc = "Cytube links go here."
		clr = 0x8FBCBB
	case "!putmetosleep":
		desc, _ = ReturnWiki()
		s.ChannelMessageSend(m.ChannelID, desc)
		return
	case "!zoltar":
		var msg string = "ZOLTAR SAYS: Make your wish, !fortune."
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	default:
		return
	}
	me := discordgo.MessageEmbed{ Title: msgtitle, Description: desc, Color: clr }
	s.ChannelMessageSendEmbed(m.ChannelID, &me)
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

	dg.AddHandler(ParseCommand)
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
