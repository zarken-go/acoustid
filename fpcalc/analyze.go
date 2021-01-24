package fpcalc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/zarken-go/acoustid"
)

type cmdRunner interface {
	CombinedOutput() ([]byte, error)
}

type analyzeErr struct {
	err error
}

func (e analyzeErr) Error() string {
	return fmt.Sprintf(`fpcalc: Analyze(%s)`, e.err.Error())
}

func AnalyzePath(path string) (acoustid.Fingerprint, error) {
	Cmd := exec.Command(findBinary(), `-json`, path)
	return parseCombinedOutput(Cmd)
}

func findBinary() string {
	EnvPath := os.Getenv(`FPCALC_PATH`)
	if EnvPath != `` {
		return EnvPath
	}
	return `fpcalc`
}

func parseCombinedOutput(Cmd cmdRunner) (acoustid.Fingerprint, error) {
	b, err := Cmd.CombinedOutput()
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok {
			if execErr.Err == exec.ErrNotFound {
				return acoustid.Fingerprint{}, analyzeErr{
					err: errors.New("executable not found"),
				}
			}
		}
		return acoustid.Fingerprint{}, analyzeErr{
			err: errors.New(string(bytes.TrimSpace(b))),
		}
	}

	var Data = struct {
		Duration    float64 `json:"duration,omitempty"`
		Fingerprint string  `json:"fingerprint,omitempty"`
	}{}

	err = json.Unmarshal(b, &Data)
	return acoustid.Fingerprint{
		Duration:    int(Data.Duration),
		Fingerprint: Data.Fingerprint,
	}, err
}
