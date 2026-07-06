// Package utils holds small shared helpers.
package utils

import "time"

// Now returns the current UTC time. Centralized so timestamps stay consistent.
func Now() time.Time { return time.Now().UTC() }
