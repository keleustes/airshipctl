package prometheus_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/prometheus"
	"opendev.org/airship/airshipctl/testutil"
)

func TestPrometheus(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "prometheus-cmd-with-defaults",
			CmdLine: "",
			Cmd:     prometheus.NewPrometheusCommand(nil),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
