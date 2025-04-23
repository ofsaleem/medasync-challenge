package medasyncchallenge

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"fmt"
)

func TestScanInput(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
	procedures := make(map[string]int)
	procedures["XXXY"] = 2
	procedures["XYXY"] = 1
	procedures["ZZZZ"] = 1
	in, _ := time.Parse(layout, "1992-01-29T19:00:00Z")
	out, _ := time.Parse(layout, "2025-04-22T00:32:00Z")
	history := History{
		in:         in,
		out:        out,
		procedures: procedures,
	}
	answer := map[string]History{
		"Omar": history,
	}

	input := `Patient Omar
Action Intake Omar 1992-01-29T19:00:00Z
Action Discharge Omar 2025-04-22T00:32:00Z
Action Treatment Omar 2000-01-01T10:00:00Z XXXY
Action Treatment Omar 2005-04-02T11:04:43Z XXXY
Action Treatment Omar 2010-09-09T09:09:09Z XYXY
Action Treatment Omar 2020-11-11T11:11:11Z ZZZZ`
	// test stdin
	reader := bufio.NewScanner(strings.NewReader(input))
	testOutput := scanInput(reader)
	if val, ok := testOutput["Omar"]; !ok {
		t.Errorf("Patient Omar not in database")
	} else {
		if val.in != answer["Omar"].in {
			t.Errorf("Omar's intake was wrong: expected %s, got %s", answer["Omar"].in, val.in)
		}
		if val.out != answer["Omar"].out {
			t.Errorf("Omar's discharge was wrong: expected %s, got %s", answer["Omar"].out, val.out)
		}
		if len(val.procedures) != len(answer["Omar"].procedures) {
			t.Errorf("Omar's procedure count was wrong: expected %d, got %d", len(answer["Omar"].procedures), len(val.procedures))
		}
	}
	// test file reading
	path := filepath.Join("testdata", "omartest.txt")
	file,_ := os.Open(path)
	reader = bufio.NewScanner(file)
	testOutput = scanInput(reader)
	if val, ok := testOutput["Omar"]; !ok {
		t.Errorf("Patient Omar not in database")
	} else {
		if val.in != answer["Omar"].in {
			t.Errorf("Omar's intake was wrong: expected %s, got %s", answer["Omar"].in, val.in)
		}
		if val.out != answer["Omar"].out {
			t.Errorf("Omar's discharge was wrong: expected %s, got %s", answer["Omar"].out, val.out)
		}
		if len(val.procedures) != len(answer["Omar"].procedures) {
			t.Errorf("Omar's procedure count was wrong: expected %d, got %d", len(answer["Omar"].procedures), len(val.procedures))
		}
	}
	file.Close()
}

func TestParseDur(t *testing.T) {
	layout := "2006-01-02T15:04:05Z"
in,_ := time.Parse(layout, "2000-01-01T10:00:00Z")
	timeSugar := func (layout string, input string) time.Time {
		out, _ := time.Parse(layout, input)
		return out
	}
	var tests = []struct {
		in, out time.Time
		want string
	}{
			{in, timeSugar(layout, "2000-01-01T11:00:00Z"), "1 hour"},
			{in, timeSugar(layout, "2000-01-01T11:01:00Z"), "1 hour and 1 minute"},
			{in, timeSugar(layout, "2000-01-01T11:01:01Z"), "1 hour, 1 minute, and 1 second"},
			{in, timeSugar(layout, "2001-01-01T10:00:00Z"), "1 year"},
			{in, timeSugar(layout, "2001-01-02T10:00:00Z"), "1 year and 1 day"},
			{in, timeSugar(layout, "2001-01-02T11:00:00Z"), "1 year, 1 day, and 1 hour"},
			{in, timeSugar(layout, "2001-01-01T11:00:01Z"), "1 year, 1 hour, and 1 second"},
			{in, timeSugar(layout, "2002-01-02T12:02:02Z"), "2 years, 2 hours, 2 minutes, and 2 seconds"},
			{in, timeSugar(layout, "2025-04-22T22:18:35Z"), "25 years, 118 days, 12 hours, 18 minutes, and 35 seconds"},
		}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s - %s", tt.in, tt.out)
		t.Run(testName, func(t *testing.T) {
			duration := tt.out.Sub(tt.in)
			gotAnswer := parseDur(duration)
			if gotAnswer != tt.want {
				t.Errorf("expected %s, got %s", tt.want, gotAnswer)
			}
		})
	}
}
