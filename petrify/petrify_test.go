package petrify

import "testing"

const (
	Input   = "(github.com/nickng/cfsm/petrify).T"
	Encoded = "LPARENgithubDOTcomSLASHnickngSLASHcfsmSLASHpetrifyRPARENDOTT"
)

func TestEncode(t *testing.T) {
	enc := Encode(Input)
	if enc != Encoded {
		t.Errorf("Encoding incorrect:\n%s\n ➔  %s\ngot %s", Input, Encoded, enc)
	}
}

func TestDecode(t *testing.T) {
	enc := Encode(Input)
	dec := Decode(enc)
	if dec != Input {
		t.Errorf("Decoding incorrect:\n%s\n ➔  %s\ngot %s", enc, Input, dec)
	}
}
