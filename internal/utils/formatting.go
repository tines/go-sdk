package utils

func SetUserAgent(ua string) string {
	if ua == "" {
		return "Tines/GoSdk"
	}
	return ua
}
