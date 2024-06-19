package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"watcherserver/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing dir to watch.")
		fmt.Println("Usage: watcherserver dir")
		return
	}

	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Join(curDir, os.Args[1])
	// check if baseDir is valid
	f, err := os.Open(baseDir)
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Println("watcherserver starting...")

	s := server.New(baseDir, 8080)
	err = s.Listen()
	if err != nil {
		log.Println(err)
	}
}
