package tmplfunc

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// QRCode encodes the provided data into a PNG-formatted image. This function takes 1, 2 or 3 arguments:
//
//	QRCode(data string)
//	QRCode(data string, width int)
//	QRCode(data string, width int, recoveryLevel int)
//
// If a parameter is not passed, the defaults are used. The default width is 512 and the default recovery level is 1 (medium).
func QRCode(data string, args ...any) ([]byte, error) {
	if len(args) > 2 {
		return nil, fmt.Errorf("can only specify 1, 2 or 3 arguments")
	}

	width := 512
	if len(args) > 0 {
		w, ok := args[0].(int)
		if !ok {
			return nil, fmt.Errorf("invalid width specified")
		}

		width = w
	}

	recovery := qrcode.Medium
	if len(args) > 1 {
		r, ok := args[1].(int)
		if !ok {
			return nil, fmt.Errorf("invalid recovery level specified")
		}

		if r < int(qrcode.Low) || r > int(qrcode.Highest) {
			return nil, fmt.Errorf("invalid recovery level specified (must be between 0 and 3)")
		}

		recovery = qrcode.RecoveryLevel(r)
	}

	by, err := qrcode.Encode(data, recovery, width)
	if err != nil {
		return nil, fmt.Errorf("qrcode: encode: %w", err)
	}

	return by, nil
}
