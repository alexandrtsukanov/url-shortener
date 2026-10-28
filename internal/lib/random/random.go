package random

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

func NewRandomString(size uint8) (string, error) {
	randUnique := rand.New(rand.NewSource(time.Now().UnixNano()))

	if size == 0 {
		return "", errors.New("alias cannot have zero length")
	}

	result := make([]string, size)

	for i := 0; i < int(size); i += 1 {
		result[i] = string(byte(randUnique.Intn(122-97+1) + 97))
	}

	return strings.Join(result, ""), nil
}
