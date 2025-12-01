package certificates

import (
	"encoding/base64"
	"strings"
)

func DecodePossibleBase64(input string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err == nil && len(decoded) > 0 && IsPEMContent(decoded) {
		return decoded
	}

	return []byte(input)
}

func IsPEMContent(data []byte) bool {
	return len(data) > 0 && (strings.Contains(string(data), "-----BEGIN") ||
		data[0] == '-')
}
