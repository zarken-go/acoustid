package acoustid

type Fingerprint struct {
	Duration    int    `json:"duration,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}
