package main

import (
	"strings"
	"unicode"
    "dist-systems/mr"
    "strconv"
)

func Map(filename string, contents string) []mr.KeyValue {

    ff := func(r rune) bool { return !unicode.IsLetter(r) }//&& !unicode.IsNumber(r)}
    words := strings.FieldsFunc(contents, ff)
    kva := []mr.KeyValue{}
    for _, w := range words {
        kv := mr.KeyValue{Key:w, Value:"1"}
        kva = append(kva, kv)
    }
    return kva
}


func Reduce(key string, values []string) string {
    return strconv.Itoa(len(values))
}


// func main() {
//     ff := func(r rune) bool { return !unicode.IsLetter(r) }//&& !unicode.IsNumber(r)}
//     // var s []string
//     // s = strings.Fields("   hello I'm santhosh 1k 23 4   ")
//     // fmt.Printf("s : %s %d", s[1], len(s))
//     // fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
//     fmt.Printf("Fields are: %q", strings.FieldsFunc(" 123 foo2 bar  baz   ", ff))
// }
