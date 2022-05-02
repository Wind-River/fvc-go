package hasher

import (
	"io"
)

type FileCollectionHasher interface {
	ReadFile(io.Reader) error
	Sum() []byte
	Hex() string
}
