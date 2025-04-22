package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type History struct {
	in         time.Time
	out        time.Time
	procedures map[string]int
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

	patients := make(map[string]History)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunks := strings.Split(scanner.Text(), " ")
		switch chunks[0] {
		case "Patient":
			m := make(map[string]int)
			patients[chunks[1]] = History{procedures: m}
		case "Action":
			switch chunks[1] {
			case "Intake":
				patient, ok := patients[chunks[2]]
				if !ok {
					fmt.Println("Error: patient", chunks[2], "not in database")
					os.Exit(1)
				}
				patient.in = parseTime(chunks[3])
				patients[chunks[2]] = patient
			case "Discharge":
				patient, ok := patients[chunks[2]]
				if !ok {
					fmt.Println("Error, patient", chunks[2], "not in database")
					os.Exit(1)
				}
				patient.out = parseTime(chunks[3])
				patients[chunks[2]] = patient
			case "Treatment":
				patient, ok := patients[chunks[2]]
				if !ok {
					fmt.Println("Error, patient", chunks[2], "not in database")
					os.Exit(1)
				}
				patient.procedures[chunks[4]] += 1
				patients[chunks[2]] = patient
			}
		}
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	// close the file when we're done
	file.Close()

	// do our output prints
	process(patients)

}

func parseTime(input string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	out, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}
	return out
}

func process(patients map[string]History) {
	for patient, history := range patients {
		procedures := len(history.procedures)
		duration := history.out.Sub(history.in)
		output := "Patient " + patient + " stayed for "
		output += parseDur(duration)
		output += fmt.Sprintf(" and had %d procedures.", procedures)
		fmt.Println(output)
	}
}

func parseDur(duration time.Duration) string {
	var years, days, hours, minutes int
	const secondsInAYear int = 60 * 60 * 24 * 365
	const secondsInADay int = 60 * 60 * 24
	const secondsInAnHour int = 60 * 60
	const secondsInAMinute int = 60
	seconds := int(duration.Seconds())
	remainder := seconds
	if seconds >= secondsInAYear {
		years = seconds / secondsInAYear
		remainder = seconds % secondsInAYear
	}
	if remainder >= secondsInADay {
		days = remainder / secondsInADay
		remainder %= secondsInADay
	}
	if remainder >= secondsInAnHour {
		hours = remainder / secondsInAnHour
		remainder %= secondsInAnHour
	}
	if remainder >= secondsInAMinute {
		minutes = remainder / secondsInAMinute
		remainder %= secondsInAMinute
	}
	output := []string{}
	if years > 0 {
		str := fmt.Sprintf("%d, ", years)
		output = append(output, str)
	}
	if days > 0 {
		str := fmt.Sprintf("%d, ", days)
		output = append(output, str)
	}
	if hours > 0 {
		str := fmt.Sprintf("%d, ", hours)
		output = append(output, str)
	}
	if minutes > 0 {
		str := fmt.Sprintf("%d, ", minutes)
		output = append(output, str)
	}
	if remainder > 0 {
		str := fmt.Sprintf("%d, ", remainder)
		output = append(output, str)
	}
	if len(output) == 0 {
		fmt.Println("exiting, patient had a 0 time duration stay")
		os.Exit(1)
	}
	if len(output) == 1 {
		return output[0]
	}
	outputStr := output[0]
	for idx := 1; idx < len(output)-2; idx++ {
		outputStr += ", " + output[idx]
	}
	outputStr += ", and " + output[len(output)-1]
	return outputStr
}
