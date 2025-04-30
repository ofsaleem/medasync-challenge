package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// to keep patient history in the map with the patient name
type History struct {
	in         time.Time
	out        time.Time
	procedures map[string]int
}

type ProcedureCount struct {
	procedure string
	count     int
}

func main() {
	// flag handling
	upgradePtr := flag.Bool("upgrade", false, "output with more, prettier, detail")
	fileNamePtr := flag.String("f", "", "file path to read")
	costNamePtr := flag.String("c", "", "cost file path to read")
	fileProvided := false
	costProvided := false
	var costMap map[string]int
	var scanner *bufio.Scanner
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "f" {
			if *fileNamePtr == "" {
				fmt.Println("Error: provide filename")
				os.Exit(1)
			}
			fileProvided = true
		}
		if f.Name == "c" {
			if *costNamePtr == "" {
				fmt.Println("Error: provide cost filename")
				os.Exit(1)
			}
			costProvided = true
		}
	})

	var costScanner *bufio.Scanner
	if costProvided {
		costFile, err := os.Open(*costNamePtr)
		if err != nil {
			panic(err)
		}
		costScanner = bufio.NewScanner(costFile)
		costMap = scanCost(costScanner)
		defer costFile.Close()
	}

	if fileProvided {
		// try to open file
		file, err := os.Open(*fileNamePtr)
		if err != nil {
			panic(err)
		}
		scanner = bufio.NewScanner(file)
		defer file.Close()
	} else {
		// if not using a file, read from STDIN
		scanner = bufio.NewScanner(os.Stdin)
	}

	// create patient map from data
	patients := scanInput(scanner)

	// do our output prints
	process(patients, costMap, *upgradePtr)
}

func scanCost(scanner *bufio.Scanner) map[string]int {
	costMap := make(map[string]int)
	for scanner.Scan() {
		chunks := strings.Split(scanner.Text(), " ")
		costMap[chunks[0]], _ = strconv.Atoi(chunks[1])
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return costMap
}

func calculateCost(procedures map[string]int, costMap map[string]int) int {
	total := 0
	for proc, times := range procedures {
		total += costMap[proc] * times
	}
	return total
}

func countProcedures(procedures map[string]int) map[string]int {
	var totalProcedures map[string]int
	for proc, times := range procedures {
		totalProcedures[proc] += times
	}
	return totalProcedures
}

func scanInput(scanner *bufio.Scanner) map[string]History {
	patients := make(map[string]History)
	// loop over each line of the input
	for scanner.Scan() {
		chunks := strings.Split(scanner.Text(), " ")
		// order of information on each line is "guaranteed"
		switch chunks[0] {
		// if we see "Patient" we know to onboard
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
	return patients
}

// basic time parse using time stdlib
func parseTime(input string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	out, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}
	return out
}

func process(patients map[string]History, costMap map[string]int, upgrade bool) {
	var totalProcedures map[string]int
	for patient, history := range patients {
		// since we used a map for treatment codes, this len() is the
		// number of unique treatments
		procedures := len(history.procedures)
		// go calculates the duration for us
		duration := history.out.Sub(history.in)
		output := "Patient " + patient + " stayed for "
		// if --upgrade flag not present, do basic calc and print
		if !upgrade {
			hours := int(duration.Hours())
			minutes := float64(duration/time.Minute) - float64(hours*60)
			// for some reason we lose seconds with above calculation so we have to re-add them
			minutes += float64(duration/time.Second-(duration/time.Minute*60)) / 60
			output += fmt.Sprintf(
				"%d.0 hours and %.1f minutes and received %d treatments,", hours, minutes, procedures,
			)
		} else {
			// if --upgrade is present, hand it over to advanced calc/print
			output += parseDur(duration)
			if procedures == 1 {
				output += fmt.Sprintf(" and received 1 treatment,")
			} else {
				output += fmt.Sprintf(" and received %d treatments,", procedures)
			}
		}
		output += fmt.Sprintf(" which cost a total of $%d.", calculateCost(history.procedures, costMap))
		fmt.Println(output)
		totalProcedures = countProcedures(history.procedures)
	}
	common := ProcedureCount{}
	for proc, times := range totalProcedures {
		if times > common.count {
			common.count = times
			common.procedure = proc
		}
	}
	output := fmt.Sprintf(" The most common procedure was %s, which was undergone %d times.", common.procedure, common.count)
	fmt.Println(output)
}

func parseDur(duration time.Duration) string {
	var years, days, hours, minutes int
	// our constants for standard lengths of time
	const secondsInAYear int = 60 * 60 * 24 * 365
	const secondsInADay int = 60 * 60 * 24
	const secondsInAnHour int = 60 * 60
	const secondsInAMinute int = 60
	seconds := int(duration.Seconds())
	// repeatedly calculate integer values and save remainders
	// so we can determine integer # of smaller time slice
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
	// this big block below adds only the time slices that
	// we actually need to print
	output := []string{}
	str := ""
	if years > 0 {
		// singular / plural logic
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
	// below block handles grammar and oxford commas
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
