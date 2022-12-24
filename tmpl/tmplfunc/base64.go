package tmplfunc

import "encoding/base64"

func Base64(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func Base64URL(in []byte) string {
	return base64.URLEncoding.EncodeToString(in)
}
