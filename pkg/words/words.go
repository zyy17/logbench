package words

import (
	"bufio"
	"bytes"
	_ "embed"
)

//go:embed words_alpha.txt
var rawWords string

var Words []string

func init() {
	buf := bytes.NewBufferString(rawWords)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		Words = append(Words, scanner.Text())
	}
}
