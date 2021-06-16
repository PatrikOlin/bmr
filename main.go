package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var path = os.ExpandEnv("$HOME") + "/.bmr"
var register = make(map[string]string)

func main() {
	readFromFile()

	markCmd := flag.NewFlagSet("mark", flag.ExitOnError)
	setMark := markCmd.String("set", "", "set mark")
	focusMark := markCmd.String("focus", "", "focus mark")

	regCmd := flag.NewFlagSet("reg", flag.ExitOnError)
	clear := regCmd.Bool("clear", false, "clear mark register")
	read := regCmd.Bool("read", false, "read register")

	if len(os.Args) < 2 {
		fmt.Println("expected 'mark' or 'reg' subcommand")
	}

	switch os.Args[1] {

	case "mark":
		markCmd.Parse(os.Args[2:])

		if *setMark != "" {
			parts := strings.Split(*setMark, ",")
			register[parts[0]] = parts[1]
			writeToFile(register)
		}

		if *focusMark != "" {
			nodeID := register[*focusMark]
			cmd := exec.Command("bspc", "node", "-f", nodeID)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				return
			}
			fmt.Println("Result: " + out.String())
		}

	case "reg":
		regCmd.Parse(os.Args[2:])

		if *read {
			dump(readFromFile())
		}

		if *clear {
			err := os.Remove(path)
			if err != nil {
				log.Fatalln("Error clearing register ", err)
			}
		}

	default:
		fmt.Println("expected 'mark' or 'reg' subcommand")
		os.Exit(1)

	}
}

func writeToFile(register map[string]string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalln("Error opening/creating file: ", err)
	}
	encoder := gob.NewEncoder(f)

	err = encoder.Encode(register)
	if err != nil {
		log.Fatalln("Error encoding map: ", err)
	}

	f.Close()
}

func readFromFile() map[string]string {
	file, err := os.Open(path)
	if err != nil {
		return make(map[string]string)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	register = make(map[string]string)
	decoder.Decode(&register)

	return register
}

func dump(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Print(string(b))
}
