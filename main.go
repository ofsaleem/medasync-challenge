package main
import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
)

type History struct {
	in time.Time
	out time.Time
	procedures int
}

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
	
	// close the file when we're done
	defer file.Close()

	patients := make(map[string]History)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunks := strings.Split(scanner.Text(), " ")
		switch chunks[0] {
			case "Patient":
				patients[chunks[0]] = History{}
			case "Action":
				switch chunks[1] {
					case "Intake":
					case "Discharge":
					case "Treatment":			
			}
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}
