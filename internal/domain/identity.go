package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func CanonicalSampleID(raw string) string {
	return strings.TrimSpace(raw)
}

func SampleFingerprint(sample Sample) string {
	payload := strings.Join([]string{
		sample.SampleID,
		sample.Source,
		string(sample.Priority),
		string(sample.Status),
		strings.Join(NormalizeTags(sample.Tags), ","),
	}, "|")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func EventFingerprint(event Event) string {
	payload := strings.Join([]string{
		event.ID,
		event.SampleID,
		string(event.Kind),
		event.Message,
		event.Actor,
		event.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000000Z"),
	}, "|")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func FingerprintPrefix(value string, length int) string {
	if length <= 0 || length >= len(value) {
		return value
	}
	return value[:length]
}

func SameIdentity(left, right Sample) bool {
	return CanonicalSampleID(left.SampleID) == CanonicalSampleID(right.SampleID)
}

func IdentityKey(sample Sample) string {
	return strings.ToLower(CanonicalSampleID(sample.SampleID))
}
