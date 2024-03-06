package utils_cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Adrephos/jeavendanc-st0263/Peer/client"
)

// Commands
var (
	DOWNLOAD = "download"
	UPLOAD   = "upload"
	SEARCH   = "search"
	PEERS    = "peers"
	LIST     = "list"
	HELP     = "help"
	EXIT     = "exit"
)

// Colors
const (
	MAGENTA = "\u001b[35m"
	YELLOW  = "\u001b[33m"
	WHITE   = "\u001b[37m"
	GREEN   = "\u001b[32m"
	BLUE    = "\u001b[34m"
	RED     = "\u001b[31m"
)

func printHelp() {
	help := `Welcome to p2p

  - peers
    Lists available peers with their URLs
    usage: peers
  - list
    Lists the files available in a peer
    usage: list <url>
  - search
    Search for a file in all available peers
    usage: search <file>
  - download
    Searches for a file and downloads it from an available peer
    usage: download <file>
  - upload
    Upload a file to a peer
    usage: upload <file> <url>
  - exit
    Exists the program
    usage: exit
  - help
    Prints this message
    usage: help`
	fmt.Println(help)
}

func logErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func splitCmd(line string) (cmd, arg string) {
	idx := strings.Index(line, " ")
	if idx == -1 {
		cmd, arg = line, ""
		return cmd, arg
	}
	cmd = line[:strings.Index(line, " ")]
	arg = line[strings.Index(line, " ")+1:]
	return cmd, arg
}

func readCommand() (string, string, error) {
	var line, cmd, arg string
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", nil
	}
	cmd, arg = splitCmd(line)
	return cmd, arg, nil
}

func resolveCommand(PClient *client.PeerClient, cmd, arg string) {
	switch cmd {
	case DOWNLOAD:
		if arg == "" { return }
		url, err := PClient.Search(arg)
		if err != nil {
			log.Println(err)
			return
		}
		PClient.Download(arg, url)
	case UPLOAD:
		if arg == "" { return }
		file, url := splitCmd(arg)
		err := PClient.Upload(file, url)
		logErr(err)
	case SEARCH:
		if arg == "" { return }
		_, err := PClient.Search(arg)
		if err != nil {
			log.Println(err)
			return
		}
	case PEERS:
		PClient.GetPeers()
	case LIST:
		if arg == "" { return }
		PClient.List(arg)
	case HELP:
		printHelp()
	case EXIT:
		return
	default:
		fmt.Printf("command %s not found\n", cmd)
		return
	}
}

func CommandLine(PClient *client.PeerClient) {
	var cmd, arg string
	var err error
	fmt.Printf("%sType \"help\" to get commands information:\n", MAGENTA)
	for cmd != EXIT {
		fmt.Printf("%s~> %s", GREEN, BLUE)
		cmd, arg, err = readCommand()
		fmt.Print(WHITE)
		if err != nil {
			continue
		}
		resolveCommand(PClient, cmd, arg)
	}
}
