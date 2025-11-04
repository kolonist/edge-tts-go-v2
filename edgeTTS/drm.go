// DRM module is used to handle DRM operations with clock skew correction.
// Currently the only DRM operation is generating the Sec-MS-GEC token value
// used in all API requests to Microsoft Edge's online text-to-speech service.

package edgeTTS

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	winEpoch      = 11644473600
	nanosInSecond = 1e9
)

// generateSecMSGEC generates the `Sec-MS-GEC` token value.
//
// This function generates a token value based on the current time in Windows file time format
// adjusted for clock skew, and rounded down to the nearest 5 minutes. The token is then hashed
// using SHA256 and returned as an uppercased hex digest.
//
// Returns:
//
//	The generated Sec-MS-GEC token value.
//
// See Also:
//
//	https://github.com/rany2/edge-tts/issues/290#issuecomment-2464956570
func generateSecMSGEC() string {
	// Get the current timestamp in Unix format with clock skew correction
	ticks := time.Now().Unix()

	// Switch to Windows file time epoch (1601-01-01 00:00:00 UTC)
	ticks += winEpoch

	// Round down to the nearest 5 minutes (300 seconds)
	ticks -= ticks % 300

	// Convert the ticks to 100-nanosecond intervals (Windows file time format)
	ticks *= nanosInSecond / 100

	// Create the string to hash by concatenating the ticks and the trusted client token
	str_to_hash := fmt.Sprintf("%d%s", ticks, trustedClientToken)

	// Compute the SHA256 hash and return the uppercased hex digest
	sum := sha256.Sum256([]byte(str_to_hash))
	return fmt.Sprintf("%X", sum)
}
