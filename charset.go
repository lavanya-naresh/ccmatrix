package main

var asciiChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// Half-width katakana (U+FF66 - U+FF9D). Each occupies a single terminal cell.
var katakanaChars = func() []rune {
	var rs []rune
	for r := rune(0xFF66); r <= 0xFF9D; r++ {
		rs = append(rs, r)
	}
	return rs
}()
