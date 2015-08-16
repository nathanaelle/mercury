package message

import (
	"time"
	"testing"
)


type MessageTest struct {
	m Message
	a string
}


var messageTest = []MessageTest{
	{ Message { "", time.Unix(0,0), "-", "-", "-", "-", "-" }, " 1970-01-01T01:00:00+01:00 - - - - -"},
	{ Message { "<0>1", time.Unix(0,0), "-", "-", "-", "-", "-" }, "<0>1 1970-01-01T01:00:00+01:00 - - - - -"},
	{ Message { "<0>1", time.Unix(0,0), "bla", "bli", "blu", "blo", "- message" }, "<0>1 1970-01-01T01:00:00+01:00 bla bli blu blo - message"},
}


func TestStringify(t *testing.T) {
	for _, tt := range messageTest {
		a := tt.m.Stringify()
		if a != tt.a {
			t.Errorf(" %v Stringify() = %v; want %v", tt.m, a, tt.a)
			continue
		}
	}
}


func Benchmark_Message_Stringify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range messageTest {
			tt.m.Stringify()
		}
	}
}
