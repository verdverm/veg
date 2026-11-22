package create

import (
	"testing"
)

func TestLooksLikeRepo(t *testing.T) {
	type tcase struct {
		str string
	  exp bool
	}
	cases := []tcase{
		{ exp: false, str: "./" },
		{ exp: false, str: "../" },
		{ exp: false, str: "../../.." },
		{ exp: false, str: "../../../foo" },
		{ exp: false, str: "foo" },
		{ exp: false, str: "foo.cue" },
		{ exp: false, str: "foo/bar" },
		{ exp: false, str: "foo/bar.cue" },
		{ exp: false, str: "foo/bar/baz.cue" },
		{ exp: true,  str: "github.com/verdverm/vegmod-cli" },
		{ exp: true,  str: "github.com/verdverm/vegmod-cli/creator" },
		{ exp: true,  str: "veg.io/hofmod-cli" },
	}

	for _, C := range cases {
		if r := looksLikeRepo(C.str); r != C.exp {
			t.Fatalf("for %q, expected %v but got %v", C.str, C.exp, r)
		}
	}
}
