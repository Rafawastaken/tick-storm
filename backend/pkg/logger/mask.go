package logger

import "strings"

// MaskEmail returns a redacted form of an email suitable for logging.
// "someone@example.com" -> "s***@e***.com"
func MaskEmail(email string) string {
	at := strings.IndexByte(email, '@')
	if at <= 0 {
		return "***"
	}
	local := email[:at]
	domain := email[at+1:]
	dot := strings.LastIndexByte(domain, '.')
	if dot <= 0 {
		return local[:1] + "***@***"
	}
	return local[:1] + "***@" + domain[:1] + "***" + domain[dot:]
}

// MaskSecret returns a short redacted form for a secret string.
// Short or empty inputs become "***"; longer ones keep a 2-char prefix/suffix.
func MaskSecret(s string) string {
	if len(s) <= 6 {
		return "***"
	}
	return s[:2] + "***" + s[len(s)-2:]
}
