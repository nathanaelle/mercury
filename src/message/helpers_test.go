package message

import (
	"strings"
	"unicode"
	"math/rand"
	"testing"
)


type FieldsTest struct {
	s string
	a []string
}


var abcd = "abcd"
var faces = "☺☻☹"
var commas = "1,2,3,4"
var dots = "1....2....3....4"


func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var makeFieldsInput = func() string {
	x := make([]byte, 1<<20)
	// Input is ~10% space, ~10% 2-byte UTF-8, rest ASCII non-space.
	for i := range x {
		switch rand.Intn(10) {
		case 0:
			x[i] = ' '
		case 1:
			if i > 0 && x[i-1] == 'x' {
				copy(x[i-1:], "χ")
				break
			}
			fallthrough
		default:
			x[i] = 'x'
		}
	}
	return string(x)
}

var fieldsInput = makeFieldsInput()


var fieldstests = []FieldsTest{
	{"", []string{}},
	{" ", []string{}},
	{" \t ", []string{}},
	{"  abc  ", []string{"abc"}},
	{"1 2 3 4", []string{"1", "2", "3", "4"}},
	{"1  2  3  4", []string{"1", "2", "3", "4"}},
	{"1\t\t2\t\t3\t4", []string{"1", "2", "3", "4"}},
	{"1\u20002\u20013\u20024", []string{"1", "2", "3", "4"}},
	{"\u2000\u2001\u2002", []string{}},
	{"\n™\t™\n", []string{"™", "™"}},
	{faces, []string{faces}},
}


var FieldsFuncTests = []FieldsTest{
	{"", []string{}},
	{"XX", []string{}},
	{"XXhiXXX", []string{"hi"}},
	{"aXXbXXXcX", []string{"a", "b", "c"}},
}





func TestFieldsFuncN(t *testing.T) {
	for _, tt := range fieldstests {
		a := FieldsFuncN(tt.s, 10, unicode.IsSpace)
		if !eq(a, tt.a) {
			t.Errorf("FieldsFunc(%q, unicode.IsSpace) = %v; want %v", tt.s, a, tt.a)
			continue
		}
	}
	pred := func(c rune) bool { return c == 'X' }
	for _, tt := range FieldsFuncTests {
		a := FieldsFuncN(tt.s, 5, pred)
		if !eq(a, tt.a) {
			t.Errorf("FieldsFunc(%q) = %v, want %v", tt.s, a, tt.a)
		}
	}

	FI_r := strings.FieldsFunc(fieldsInput, unicode.IsSpace)
	FI_t := FieldsFuncN(fieldsInput, len(fieldsInput)/20, unicode.IsSpace)

	if !eq(FI_r, FI_t) {
		t.Errorf("FieldsFunc(FI) = %v, want %v", len(FI_t), len(FI_t))
	}
}



/*

func Benchmark_strings_FieldsFunc(b *testing.B) {
	b.SetBytes(int64(len(fieldsInput)))
	for i := 0; i < b.N; i++ {
		strings.FieldsFunc(fieldsInput, unicode.IsSpace)
	}
}


func Benchmark_FieldsFuncN(b *testing.B) {
	b.SetBytes(int64(len(fieldsInput)))
	for i := 0; i < b.N; i++ {
		FieldsFuncN(fieldsInput, len(fieldsInput)/20, unicode.IsSpace)
	}
}

*/
