# MedaSync Engineering Challenge
This repository represents my work to complete the MedaSync engineering challenge. I've chosen to write my program in Go, with an additional Bash script used for certain tests.

## Usage
Use the included binary `medasync-challenge` to run the program, or `go run main.go` in place will also work.

The program supports a couple flags:

`-f filename` to read input from a provided file

`--upgrade` to give more granular time calculations (e.g. years, days) and provide slightly prettier output.

`./medasync-challenge -f example_patients.txt --upgrade`

`go run main.go -f example_patients.txt --upgrade`

Be sure to enable executable permissions on the binary.

Both flags are optional. If `-f` is not provided, the program will read from Standard Input, so you can pipe to it if you'd like:

`cat example_patients.txt | ./medasync-challenge --upgrade`
