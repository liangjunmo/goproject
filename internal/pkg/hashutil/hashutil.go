package hashutil

import (
	"crypto/sha1"
)

func Sha1ByteToByte(s []byte) []byte {
	hash := sha1.New()
	hash.Write(s)
	return hash.Sum(nil)
}

func Sha1StringToString(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return string(hash.Sum(nil))
}

func Sha1ByteToString(s []byte) string {
	hash := sha1.New()
	hash.Write(s)
	return string(hash.Sum(nil))
}

func Sha1StringToByte(s string) []byte {
	hash := sha1.New()
	hash.Write([]byte(s))
	return hash.Sum(nil)
}
