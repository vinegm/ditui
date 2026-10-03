package report

import (
	"fmt"
	"slices"
	"unicode"

	"github.com/vinegm/ditui/internal/morse"
)

// Bucket is one label of a breakdown, like "Q −−•−" or "4 elements".
type Bucket struct {
	Label   string
	Correct int
	Asked   int
	Median  float64
}

// Breakdown groups the attempts along one dimension of the character.
type Breakdown struct {
	Name    string
	Buckets []Bucket // in the order each label was first asked
}

// Breakdowns are the character, code length and kind profiles of a session.
func Breakdowns(attempts []morse.Attempt) []Breakdown {
	return []Breakdown{
		{"character", group(attempts, CharLabel)},
		{"length", group(attempts, lengthLabel)},
		{"kind", group(attempts, kindLabel)},
	}
}

// CharLabel names an attempt's character along with its code: "L •−••".
func CharLabel(attempt morse.Attempt) string {
	return string(attempt.Char) + " " + morse.Pretty(attempt.Code())
}

func lengthLabel(attempt morse.Attempt) string {
	if elements := len(attempt.Code()); elements != 1 {
		return fmt.Sprintf("%d elements", elements)
	}

	return "1 element"
}

func kindLabel(attempt morse.Attempt) string {
	switch {
	case unicode.IsLetter(attempt.Char):
		return "letter"

	case unicode.IsDigit(attempt.Char):
		return "digit"

	default:
		return "punctuation"
	}
}

// Worst is the label with the lowest accuracy (slowest breaks ties), and false
// if every bucket was answered perfectly.
func (breakdown Breakdown) Worst() (string, bool) {
	if len(breakdown.Buckets) == 0 {
		return "", false
	}

	worst := breakdown.Buckets[0]

	for _, bucket := range breakdown.Buckets[1:] {
		if compareWorst(bucket, worst) < 0 {
			worst = bucket
		}
	}

	return worst.Label, worst.Correct < worst.Asked
}

// Ranked is the buckets sorted worst-first: lowest accuracy, then slowest.
func (breakdown Breakdown) Ranked() []Bucket {
	ranked := slices.Clone(breakdown.Buckets)
	slices.SortStableFunc(ranked, compareWorst)

	return ranked
}

// compareWorst orders a before b when a has the lower accuracy, or the same
// accuracy and a slower median.
func compareWorst(a, b Bucket) int {
	// Cross-multiplied so equal ratios compare equal without float division.
	if left, right := a.Correct*b.Asked, b.Correct*a.Asked; left != right {
		if left < right {
			return -1
		}

		return 1
	}

	switch {
	case a.Median > b.Median:
		return -1
	case a.Median < b.Median:
		return 1
	default:
		return 0
	}
}

func group(attempts []morse.Attempt, key func(morse.Attempt) string) []Bucket {
	var labels []string
	rows := map[string][]morse.Attempt{}

	for _, attempt := range attempts {
		label := key(attempt)

		if _, seen := rows[label]; !seen {
			labels = append(labels, label)
		}

		rows[label] = append(rows[label], attempt)
	}

	buckets := make([]Bucket, len(labels))

	for i, label := range labels {
		buckets[i] = Bucket{
			Label:   label,
			Correct: CountCorrect(rows[label]),
			Asked:   len(rows[label]),
			Median:  Median(rows[label]),
		}
	}

	return buckets
}
