package calico_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/calico"
	"opendev.org/airship/airshipctl/testutil"
)

func TestCalico(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "calico-cmd-with-defaults",
			CmdLine: "",
			Cmd:     calico.NewCalicoCommand(nil),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
