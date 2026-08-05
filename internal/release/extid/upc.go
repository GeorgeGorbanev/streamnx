package extid

import "strings"

func NormalizeUPC(rawUPC string) (string, bool) {
	upcValue := strings.TrimSpace(rawUPC)
	if len(upcValue) != 12 && len(upcValue) != 13 {
		return "", false
	}
	if !upcValidDigits(upcValue) || !upcValidCheckDigit(upcValue) {
		return "", false
	}
	if len(upcValue) == 13 && upcValue[0] == '0' {
		return upcValue[1:], true
	}
	return upcValue, true
}

func upcValidDigits(upcValue string) bool {
	for upcIndex := range upcValue {
		if upcValue[upcIndex] < '0' || upcValue[upcIndex] > '9' {
			return false
		}
	}
	return true
}

func upcValidCheckDigit(upcValue string) bool {
	upcTotal := 0
	for upcIndex := range upcValue {
		upcDigit := int(upcValue[upcIndex] - '0')
		if (len(upcValue)-upcIndex)%2 == 0 {
			upcDigit *= 3
		}
		upcTotal += upcDigit
	}
	return upcTotal%10 == 0
}
