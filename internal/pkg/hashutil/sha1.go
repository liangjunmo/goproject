package hashutil

import (
	"crypto/sha1"
)

func SHA1ByteToByte(s []byte) []byte {
	hash := sha1.New()
	hash.Write(s)
	return hash.Sum(nil)
}

func SHA1ByteToString(s []byte) string {
	hash := sha1.New()
	hash.Write(s)
	return string(hash.Sum(nil))
}

func SHA1StringToString(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return string(hash.Sum(nil))
}

func SHA1StringToByte(s string) []byte {
	hash := sha1.New()
	hash.Write([]byte(s))
	return hash.Sum(nil)
}
