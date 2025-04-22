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
	arguments := len(os.Args)
	upgrade := false
	// must provide filename on command line
	if arguments < 2 || (arguments == 2 && os.Args[1] == "--upgrade-output") {
		fmt.Println("Error: provide filename")
		os.Exit(1)
	}
	if arguments == 3 {
		if os.Args[1] == "--upgrade-output" {
			upgrade = true
		}
	}
	// try to open file
	file, err := os.Open(os.Args[arguments-1])
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
	if upgrade {
		processUpgrade(patients)
	} else {
		process(patients)
	}

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
		hours := int(duration.Hours())
		minutes := float64(duration / time.Minute) - float64(hours * 60)
		output := fmt.Sprintf("Patient %s stayed for %d.0 hours and %.1f minutes and received %d treatments", patient, hours, minutes, procedures)
		fmt.Println(output)
	}
}

func processUpgrade(patients map[string]History) {
	for patient, history := range patients {
		procedures := len(history.procedures)
		duration := history.out.Sub(history.in)
		output := "Patient " + patient + " stayed for "
		output += parseDur(duration)
		if procedures == 1 {
			output += fmt.Sprintf(" and had 1 procedure.")
		} else {
			output += fmt.Sprintf(" and had %d procedures.", procedures)
		}
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
	str := ""
	if years > 0 {
		if years == 1 {
			str = fmt.Sprintf("1 year")
		} else {
			str = fmt.Sprintf("%d years", years)
		}
		output = append(output, str)
	}
	if days > 0 {
		if days == 1 {
			str = fmt.Sprintf("1 day")
		} else {
			str = fmt.Sprintf("%d days", days)
		}
		output = append(output, str)
	}
	if hours > 0 {
		if hours == 1 {
			str = fmt.Sprintf("1 hour")
		} else {
			str = fmt.Sprintf("%d hours", hours)
		}
		output = append(output, str)
	}
	if minutes > 0 {
		if minutes == 1 {
			str = fmt.Sprintf("1 minute")
		} else {
			str = fmt.Sprintf("%d minutes", minutes)
		}
		output = append(output, str)
	}
	if remainder > 0 {
		if remainder == 1 {
			str = fmt.Sprintf("1 second")
		} else {
			str = fmt.Sprintf("%d seconds", remainder)
		}
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
	for idx := 1; idx <= len(output)-2; idx++ {
		outputStr += ", " + output[idx]
	}
	if len(output) > 2 {
		outputStr += ","
	}
	outputStr += " and " + output[len(output)-1]
	return outputStr
}
