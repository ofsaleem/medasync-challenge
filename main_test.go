package main

import (
	"bufio"
	"strings"
	"testing"
	"time"
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
}
