package translit

import "strings"

var (
	singleLatinToCyrillicReplacer = strings.NewReplacer(
		"a", "а", "b", "б", "c", "ц", "d", "д", "e", "е", "f", "ф", "g", "г",
		"h", "х", "i", "и", "j", "ж", "k", "к", "l", "л", "m", "м", "n", "н",
		"o", "о", "p", "п", "q", "к", "r", "р", "s", "с", "t", "т", "u", "у",
		"v", "в", "w", "в", "x", "кс", "y", "ы", "z", "з",
		"A", "А", "B", "Б", "C", "Ц", "D", "Д", "E", "Е", "F", "Ф", "G", "Г",
		"H", "Х", "I", "И", "J", "Ж", "K", "К", "L", "Л", "M", "М", "N", "Н",
		"O", "О", "P", "П", "Q", "К", "R", "Р", "S", "С", "T", "Т", "U", "У",
		"V", "В", "W", "В", "X", "КС", "Y", "Ы", "Z", "З",
	)
	multipleLatinToCyrillicReplacer = strings.NewReplacer(
		"shch", "щ", "sh", "ш", "ch", "ч", "zh", "ж", "yo", "ё",
		"Shch", "Щ", "Sh", "Ш", "Ch", "ч", "Zh", "Ж", "Yo", "ё",
	)
	cyryllicToLatinReplacer = strings.NewReplacer(
		"а", "a", "б", "b", "в", "v", "г", "g", "д", "d", "е", "e", "ё", "yo",
		"ж", "zh", "з", "z", "и", "i", "й", "i", "к", "k", "л", "l", "м", "m",
		"н", "n", "о", "o", "п", "p", "р", "r", "с", "s", "т", "t", "у", "u",
		"ф", "f", "х", "h", "ц", "c", "ч", "ch", "ш", "sh", "щ", "shch",
		"ъ", "", "ы", "y", "ь", "", "э", "e", "ю", "yu", "я", "ya",
		"А", "A", "Б", "B", "В", "V", "Г", "G", "Д", "D", "Е", "E", "Ё", "Yo",
		"Ж", "Zh", "З", "Z", "И", "I", "Й", "I", "К", "K", "Л", "L", "М", "M",
		"Н", "N", "О", "O", "П", "P", "Р", "R", "С", "S", "Т", "T", "У", "U",
		"Ф", "F", "Х", "H", "Ц", "C", "Ч", "Ch", "Ш", "Sh", "Щ", "Shch",
		"Ъ", "", "Ы", "Y", "Ь", "", "Э", "E", "Ю", "Yu", "Я", "Ya",
	)
)

func CyrillicToLatin(input string) string {
	return cyryllicToLatinReplacer.Replace(input)
}

func LatinToCyrillic(input string) string {
	return singleLatinToCyrillicReplacer.Replace(
		multipleLatinToCyrillicReplacer.Replace(input),
	)
}
