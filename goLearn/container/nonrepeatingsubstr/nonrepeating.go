package main

func lengthOfNonRepeatingSubStr(s string) int {
	lastOccured := make(map[rune]int)
	start := 0
	maxLength := 0

	for i, ch := range []rune(s) {
		if lastI, ok := lastOccured[ch]; ok && lastI >= start {
			start = lastI + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccured[ch] = i
	}
	return maxLength
}
