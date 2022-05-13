package legacy

import (
	"bytes"

	"github.com/pkg/errors"
)

// Upgrade converts a FVC0 to FVC1, since they are the same payload format, just with a header added
func Upgrade(fvc0 []byte) ([]byte, error) {
	if !Valid(fvc0) {
		// Check if already v1
		if len(fvc0) == 25 && bytes.Compare(fvc0[0:5], []byte("FVC1\000")) == 0 {
			return fvc0, nil
		}

		return nil, errors.New("not a valid V0 or V1 File Verification Code")
	}

	return append([]byte("FVC1\000"), fvc0...), nil
}
