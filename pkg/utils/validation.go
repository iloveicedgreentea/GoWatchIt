package utils

// validate scheme

func ValidateHttpScheme(scheme string) bool {
	return scheme == "http://" || scheme == "https://"
}