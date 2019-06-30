package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

var highScore, score int32

func saveHighscore(hs int32) {
	f, err := os.Create(getScorePath())
	check(err)

	sEnc := base64.StdEncoding.EncodeToString([]byte(formatInt32(hs)))

	_, err = f.WriteString(sEnc)
	f.Sync()
}

func readHighscore() int32 {
	var f *os.File
	var err error
	p := getScorePath()

	if _, ferr := os.Stat(p); os.IsNotExist(ferr) {
		f, err = os.Create(p)
		check(err)
	} else {
		f, err = os.Open(p)
		check(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil || fi.Size() <= 0 {
		fmt.Printf("Failed to read score file status %s\n", err)
		return 0
	}

	b := make([]byte, fi.Size())
	_, err = f.Read(b)
	check(err)

	sDec, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		fmt.Println("decode error:", err)
		return 0
	}

	r, _ := strconv.Atoi(string(sDec))
	return int32(r)
}

func getScorePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	p := filepath.FromSlash(fmt.Sprintf("%s/.gosnake", usr.HomeDir))
	return p
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
