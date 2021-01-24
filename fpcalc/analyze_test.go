package fpcalc

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarken-go/acoustid"
)

type fakeCmdOutput struct {
	Data []byte
	Err  error
}

func (f fakeCmdOutput) CombinedOutput() ([]byte, error) {
	return f.Data, f.Err
}

func TestFindBinary(t *testing.T) {
	assert.Equal(t, `fpcalc`, findBinary())

	_ = os.Setenv(`FPCALC_PATH`, `/full/path/fpcalc`)
	t.Cleanup(func() {
		_ = os.Unsetenv(`FPCALC_PATH`)
	})

	assert.Equal(t, `/full/path/fpcalc`, findBinary())
}

func TestParseOutputSuccess(t *testing.T) {
	Resp, err := parseCombinedOutput(&fakeCmdOutput{
		Data: []byte(`{"duration":10.5,"fingerprint":"abcd1234"}`),
	})
	assert.Nil(t, err)
	assert.Equal(t, 10, Resp.Duration)
	assert.Equal(t, `abcd1234`, Resp.Fingerprint)
}

func TestParseOutputErr(t *testing.T) {
	Resp, err := parseCombinedOutput(&fakeCmdOutput{
		Data: []byte(`ERROR: Could not open the input file (No such file or directory)`),
		Err:  &exec.Error{},
	})
	assert.IsType(t, analyzeErr{}, err)
	assert.EqualError(t, err, `fpcalc: Analyze(ERROR: Could not open the input file (No such file or directory))`)
	assert.Equal(t, acoustid.Fingerprint{}, Resp)
}

func TestParseOutputBinaryNotFound(t *testing.T) {
	Resp, err := parseCombinedOutput(&fakeCmdOutput{
		Err: &exec.Error{
			Err: exec.ErrNotFound,
		},
	})
	assert.IsType(t, analyzeErr{}, err)
	assert.EqualError(t, err, `fpcalc: Analyze(executable not found)`)
	assert.Equal(t, acoustid.Fingerprint{}, Resp)
}
