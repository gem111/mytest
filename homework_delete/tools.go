package homework_delete

import "unicode"

func UnderscoreName(str string) string {
	var buf []byte
	for i, v := range str {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}
