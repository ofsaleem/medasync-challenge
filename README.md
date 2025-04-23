# MedaSync Engineering Challenge
This repository represents my work to complete the MedaSync engineering challenge. I've chosen to write my program in Go, with an additional Bash script used for certain tests.

## Usage
Use the included binary `medasync-challenge` to run the program, or `go run main.go` in its place will also work.

### Flags
The program supports a couple flags:

`-h` to see flag usage

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

The program does not support reading multiple files with the `-f` flag, but you can pipe multiple in at once.

### Testing
To run the unit tests I wrote (in `main_test.go`) run `go test` .

To run the end to end (from file parsing to output) tests, make sure `e2e.sh` has execution permissions and run it with `./e2e.sh` . This script generates program output based on input files living in the `testdata` folder, and `diff`s it with the answer files I have also living there. You can add your own if you like, incrementally named matching `testN.txt` and `answersN.txt` can be put in the `testdata` folder and the script will automatically run them. Be careful of extra newlines and trailing spaces in your answers files if you do this, otherwise the test will fail and the output will be hard to discern (can't highlight a space, after all).

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
#### --upgrade
When the `--upgrade` flag is present, the output is handled by a different function that does more calculations and handles grammar as well. This function supports years, days, and seconds as well as the original hours and minutes required. It distinguishes between singular and plural values (`1 treatment` and `1 hour` v. `1 treatments` and `1 hours`) and will not print 0 values (`1 year and 1 hour` rather than `1 year, 0 days, and 1 hour`). It also joins clauses following English grammar rules, including providing an oxford comma. This function also adds a period at the end of the sentence. Fun!

## Assumptions
Most important assumptions I make are those provided in the problem statement
* A correct input will always have a patient before any actions taken on that patient
* Every patient should have both an Intake and Discharge
* The Discharge date and time should always be after the Intake date and time
* Every patient should have at least one treatment

Others:
* All input provided will be correct, or behavior of the program is not guaranteed
* Patients have more-or-less human lifespans
    * Go's `time` library only supports durations of about 290 years
* I can ignore input that doesn't match expected format. So if a line doesn't start with `Patient` or `Action`, or the `Action` isn't one of the ones we expect, the program ignores the line completely.
* Data on the same line is separated by spaces and not tabs
* No different patients have the same name
* Timestamps on `Treatment` lines can be ignored, since the problem statement doesn't mention tracking these at all.

## Limitations / Room for Improvement
* I could leverage a custom struct and maybe some more `time` functions for the upgraded calculations to make that section of the code prettier or shorter, but I opted not to since that section is out of scope of the problem anyway
* Since I'm making assumptions on input being correct, error handling exists mostly to panic and quit if something is wrong. For example if a patient is not in the map when I attempt to perform an `Action` the program errors and quits. I could print these errors, maybe to a logfile, and let the program continue to process other lines.
* Similarly I don't have checks in place in case the duration is negative, which should only happen if the `Intake` time is after the `Discharge` time.
* I don't have particular logic to handle leap years or seconds. For leap years, the actual amount of time printed will be correct, but it may not be grouped correctly. For example, any span of 4 years will include 1 leap year. But a duration of exactly 4 calendar years would print `4 years and 1 day` rather than `4 years`. Leap seconds are not accounted for at all due to a Go limitation. Sorry!
* The `e2e.sh` script does not delete the incrementally named `outN.txt` files it generates during testing.
* There's no safety around input size, so providing a large enough input could definitely crash the program.
* The program only supports one file being provided at a time, if you want to read multiple files you could `cat` them all and pipe them into the binary.
* Data provided for patients with the same name will overwrite. So if you issue `Intake` twice for John, it will replace the first timestamp with the new one. This means you can't have two different patients named John and expect to get separate output for them. This does mean however there's an amount of idempotency provided by the program; if your input file accidentally includes the **same** data multiple times you will stil get accurate outputs
* Program does not rely on external memory at all so each time you run the program it's fresh.
