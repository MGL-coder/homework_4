package main

import (
	"github.com/MGL-coder/homework_4/tetris"
	"log"
)

func main() {
	err := tetris.Tetris("file.txt")
	if err != nil {
		log.Fatalln(err)
	}
}
