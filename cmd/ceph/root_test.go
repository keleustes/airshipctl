package ceph_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/ceph"
	"opendev.org/airship/airshipctl/testutil"
)

func TestCeph(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "ceph-cmd-with-defaults",
			CmdLine: "",
			Cmd:     ceph.NewCephCommand(nil),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
