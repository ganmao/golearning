package main

import (
	"SMP/mlib"
	"SMP/mp"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var lib *library.MusicManager
var id int = 1
var ctrl, signal chan int

func handleLibCommands(tokens []string) {
	switch tokens[1] {
	case "list":
		for i := 0; i < lib.Len(); i++ {
			e, _ := lib.Get(i)
			fmt.Println(i+1, ":", e.Name, e.Artist, e.Source, e.Type)
		}
	case "add":
		{
			if len(tokens) == 7 {
				id++
				lib.Add(&library.MusicEntry{strconv.Itoa(id),
					tokens[2], tokens[3], tokens[4], tokens[5], tokens[6]})
			} else {
				fmt.Println("USAGE: lib add <name><artist><sect><source><type>")
			}
		}
	case "del":
		if len(tokens) == 3 {
			n := lib.RemoveByName(tokens[2])
			fmt.Println("remove music numbers = ", n)
		} else {
			fmt.Println("USAGE: lib del <name>")
		}
	default:
		fmt.Println("Unrecognized lib command:", tokens[1])
	}
}

func handlePlayCommand(tokens []string) {
	ch_data := make(chan float64, 4)

	if len(tokens) != 2 {
		fmt.Println("USAGE: play <name>")
		return
	}
	e := lib.Find(tokens[1])
	if e == nil {
		fmt.Println("The music", tokens[1], "does not exist.")
		return
	}
	//go mp.Play(e.Source, e.Type)
	go mp.Play(e.Source, e.Type, ch_data)

	if t, ok := <-ch_data; ok {
		fmt.Printf("%s play %v seconds\n", e.Source, t)
	}
	defer close(ch_data)
}

func main() {
	fmt.Println(`
Enter following commands to control the player:
lib list -- View the existing music lib
lib add <name><artist><sect><source><type> -- Add a music to the music lib
lib del <name> -- Remove the specified music from the lib
play <name> -- Play the specified music
`)

	lib = library.NewMusicManager()
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command-> ")
		rawLine, _, _ := r.ReadLine()
		line := string(rawLine)
		if line == "q" || line == "e" {
			break
		}
		tokens := strings.Split(line, " ")
		//fmt.Println("tokens:", len(tokens))
		if tokens[0] == "lib" {
			handleLibCommands(tokens)
		} else if tokens[0] == "play" {
			handlePlayCommand(tokens)
		} else {
			fmt.Println("Unrecognized command:", tokens[0])
		}
	}
}
