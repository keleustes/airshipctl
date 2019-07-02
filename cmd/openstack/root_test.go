package openstack_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/openstack"
	"opendev.org/airship/airshipctl/testutil"
)

func TestOpenStack(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "openstack-cmd-with-defaults",
			CmdLine: "",
			Cmd:     openstack.NewOpenStackCommand(nil),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
