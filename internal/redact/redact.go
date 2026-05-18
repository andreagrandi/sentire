package redact

import "strings"

// Placeholder is the string used to replace a redacted secret.
const Placeholder = "[REDACTED]"

// Secret returns s with every occurrence of secret replaced by Placeholder.
// An empty secret is treated as a no-op so callers can apply this even when
// the token has not been loaded yet.
func Secret(s, secret string) string {
	if secret == "" {
		return s
	}
	return strings.ReplaceAll(s, secret, Placeholder)
}
