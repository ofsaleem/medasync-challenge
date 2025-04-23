# MedaSync Engineering Challenge
This repository represents my work to complete the MedaSync engineering challenge. I've chosen to write my program in Go, with an additional Bash script used for certain tests.

## Usage
Use the included binary `medasync-challenge` to run the program, or `go run main.go` in its place will also work.

### Flags
The program supports a couple flags:

`-f filename` to read input from a provided file

`--upgrade` to give more granular time calculations (e.g. years, days) and provide slightly prettier output.

```
./medasync-challenge -f example_patients.txt --upgrade
```

```
go run main.go -f example_patients.txt --upgrade
```

Be sure to enable executable permissions on the binary.

Both flags are optional. If `-f` is not provided, the program will read from Standard Input, so you can pipe to it if you'd like:

```
cat example_patients.txt | ./medasync-challenge --upgrade
```

The program does not support multiple file inputs.

## Design
For programming convenience and performance I opted to store patient data in a `map` with the patient name as the key and a custom `History` struct as the value. This lets me quickly and consistently perform writes and lookups to patient data, regardless of the order I receive the information. The patient's treatment history is also stored in a `map` so that I can quickly count only unique treatments.

The program parses the command-line flags and then immediately moves to parsing input (via file or via `STDIN`). My approach to input parsing was to handle the input line-by-line. This is in large part due to assumptions and examples provided in the problem statement document, showing that each line independently had all the information needed to parse each action. 

The other major design consideration taken here was to only ingest information during file parsing, and hold all calculations and output until after I was done with the file. The main reason to do it this way was that only the `Patient` directive had any guarantees on order. Doing processing at the end means the program does not care if `Discharge` is provided before `Intake`, or if a patient's `Treatment` lines are mixed in with another's. To support this, I stored the timestamps for intake and discharge in a patient `History` struct and calculated the duration between them during processing, at the end.

### Output
Ouput only prints to STDOUT but it can of course be piped to a file via the command-line if desired
```
./medasync-challenge -f example_patients.txt > output.txt
```
By default the program tries to match the output style in the problem statement:
```
Patient John stayed for 222.0 hours and 13.0 minutes and received 4 treatments
```
I add the `.0` to the hours statically in the print statement, because it will always be an integer result since we break out the minutes. I made an assumption that we only wanted to preserve one decimal place for the minutes, so there will be some rounding involved if there are any stays that involve seconds. These results are achieved using Go's `time` standard library, which calculates a duration for me and I do some fun math to get the remainders in minutes and then in seconds.
