package randutil

import (
	"math/rand"
	"strings"
	"time"
)

type RandomStrItem []string

var (
	RandomStrUpper RandomStrItem = []string{
		"A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
	}
	RandomStrLower RandomStrItem = []string{
		"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z",
	}
	RandomStrNumber RandomStrItem = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
	}
)

func RandomStr(length int, items ...RandomStrItem) string {
	strs := make([]string, 0)
	if len(items) == 0 {
		strs = append(strs, RandomStrUpper...)
		strs = append(strs, RandomStrLower...)
		strs = append(strs, RandomStrNumber...)
	} else {
		for _, item := range items {
			strs = append(strs, item...)
		}
	}
	rand.Seed(time.Now().Unix())
	str := strings.Builder{}
	for i := 0; i < length; i++ {
		index := rand.Intn(len(strs))
		str.WriteString(strs[index])
	}
	return str.String()
}
