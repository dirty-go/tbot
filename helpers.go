package tbot

import "strconv"

// itoa renders an int for a form field.
func itoa(n int) string { return strconv.Itoa(n) }

// itoa64 renders an int64 for a form field.
func itoa64(n int64) string { return strconv.FormatInt(n, 10) }

// formatFloat renders a float64 for a form field without trailing zeros.
func formatFloat(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
