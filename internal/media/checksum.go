package media

import (
	"crypto/sha256"
	"encoding/hex"
)

func buildChecksum(bytes []byte) string {
	sha := sha256.Sum256(bytes)
	return hex.EncodeToString(sha[:])
}
