package main
import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	// must provide filename on command line
	if len(os.Args) < 2 {
		fmt.Println("Error: provide filename")
		os.Exit(1)
	}
	// try to open file
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	

}
