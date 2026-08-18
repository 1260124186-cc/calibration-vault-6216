package domain

import "time"

type Release struct {
	Operator    string    `json:"operator"`
	Destination string    `json:"destination"`
	ReleasedAt  time.Time `json:"released_at"`
}

type ReleaseInput struct {
	Operator    string `json:"operator"`
	Destination string `json:"destination"`
}

func NewRelease(input ReleaseInput, now time.Time) Release {
	return Release{
		Operator:    input.Operator,
		Destination: input.Destination,
		ReleasedAt:  now,
	}
}

func (r Release) Clone() Release {
	return Release{
		Operator:    r.Operator,
		Destination: r.Destination,
		ReleasedAt:  r.ReleasedAt,
	}
}

func (r Release) IsFor(destination string) bool {
	return r.Destination == destination
}

func SupportedDestinations() []string {
	return []string{"lab-east", "lab-west", "lab-central"}
}
