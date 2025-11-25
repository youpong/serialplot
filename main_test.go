package main

import (
	"bufio"
	"regexp"
	"testing"
)

func TestMockReader(t *testing.T) {
	re := regexp.MustCompile(`^A[01]:\d+$`)

	m := &MockReader{}
	s := bufio.NewScanner(m)

	if !s.Scan() {
		t.Errorf("")
	}
	//	fmt.Printf("DEBUG: %s\n", s.Text())
	//	fmt.Printf("DEBUG: %s\n", s.Text())
	if !re.MatchString(s.Text()) {
		t.Errorf("")
	}
	// s.Scan()
	// fmt.Printf("DEBUG: %s\n", s.Text())
}
