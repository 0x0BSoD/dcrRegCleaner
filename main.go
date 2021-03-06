package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	confPath string

	listCommand *flag.FlagSet

	pruneCommand *flag.FlagSet
	pruneCheck   bool
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("usage: regCleaner <command> [<args>]")
		fmt.Println("	list	list registry content")
		fmt.Println("	prune	remove old images, keep configured number of items")
		os.Exit(2)
	}

	flag.StringVar(&confPath, "config", "config.yml", "path to config in YML format")

	listCommand = flag.NewFlagSet("list", flag.ExitOnError)

	pruneCommand = flag.NewFlagSet("prune", flag.ExitOnError)
	pruneCommand.BoolVar(&pruneCheck, "check", false, "just list images to remove")

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "prune":
		pruneCommand.Parse(os.Args[2:])
	default:
		log.Printf("%s is not valid command \n", os.Args[1])
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()

	config, err := parseConfig(confPath)
	if err != nil {
		log.Printf("cant read config %s", confPath)
		os.Exit(-1)
	}
	client := NewClient(config)

	if listCommand.Parsed() {
		err = printRegContent(client)
		if err != nil {
			log.Println("command failed, err:", err)
			os.Exit(-1)
		}
	}

	if pruneCommand.Parsed() {
		client.RmCheck = pruneCheck
		err = removeOldImages(client)
		if err != nil {
			log.Println("command failed, err:", err)
			os.Exit(-1)
		}
	}
}
