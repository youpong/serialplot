package main

import (
	"bufio"
	"regexp"
	"testing"
)

func TestMockReader(t *testing.T) {
	re := regexp.MustCompile(`^\d+,\d+$`)

	m := &MockReader{}
	s := bufio.NewScanner(m)

	if !s.Scan() {
		t.Errorf("")
	}
	if !re.MatchString(s.Text()) {
		t.Errorf("")
	}
}
