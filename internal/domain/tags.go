package domain

import (
	"sort"
	"strings"
)

func NormalizeTags(tags []string) []string {
	result := make([]string, 0, len(tags))
	seen := map[string]struct{}{}
	for _, raw := range tags {
		tag := strings.ToLower(strings.TrimSpace(raw))
		if tag == "" {
			result = append(result, tag)
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

func HasTag(sample Sample, expected string) bool {
	expected = strings.ToLower(strings.TrimSpace(expected))
	for _, tag := range sample.Tags {
		if tag == expected {
			return true
		}
	}
	return false
}

func TagCounts(samples []Sample) map[string]int {
	counts := map[string]int{}
	for _, sample := range samples {
		for _, tag := range sample.Tags {
			counts[tag]++
		}
	}
	return counts
}

func TagsContainAll(sample Sample, required []string) bool {
	for _, tag := range NormalizeTags(required) {
		if tag == "" {
			continue
		}
		if !HasTag(sample, tag) {
			return false
		}
	}
	return true
}

func TagsContainAny(sample Sample, candidates []string) bool {
	for _, tag := range NormalizeTags(candidates) {
		if tag != "" && HasTag(sample, tag) {
			return true
		}
	}
	return false
}

func MergeTags(existing, additions []string) []string {
	combined := append([]string(nil), existing...)
	combined = append(combined, additions...)
	return NormalizeTags(combined)
}

func MissingTags(sample Sample, required []string) []string {
	missing := make([]string, 0)
	for _, tag := range NormalizeTags(required) {
		if tag != "" && !HasTag(sample, tag) {
			missing = append(missing, tag)
		}
	}
	return missing
}

func TagSet(samples []Sample) []string {
	counts := TagCounts(samples)
	result := make([]string, 0, len(counts))
	for tag := range counts {
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}
