package image

import "fmt"

func CheckMimeType(mimeType string) (string, error) {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	default:
		return "", fmt.Errorf("unsupported media type: %v", mimeType)
	}
}
