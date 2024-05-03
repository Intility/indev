// Package pipeline collects some light telemetry to be able to improve icpctl over time.
// We're aware how important privacy is and value it ourselves, so we have
// the following rules:
// 1. We only collect anonymized data – nothing that is personally identifiable
// 2. Data is only stored in ISAE 3000 (SOC 2 equivalent) compliant systems, and we are ISAE 3000 compliant ourselves.
// 3. Users should always have the ability to opt-out.
package pipeline

import (
	"runtime/trace"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/intility/icpctl/internal/telemetry"
)

func Telemetry() *TelemetryMiddleware {
	return &TelemetryMiddleware{}
}

type TelemetryMiddleware struct{}

// TelemetryMiddleware implements interface Middleware (compile-time check).
var _ Middleware = (*TelemetryMiddleware)(nil)

func (m *TelemetryMiddleware) Handle(cmd *cobra.Command, args []string, next NextFunc) error {
	telemetry.Start()

	runErr := next(cmd, args)

	defer trace.StartRegion(cmd.Context(), "telemetryPostRun").End()
	defer telemetry.Stop()

	meta := telemetry.Metadata{} //nolint:exhaustruct

	subCmd, flags, err := getSubcommand(cmd, args)
	if err != nil {
		// Ignore invalid commands/flags.
		return runErr
	}

	meta.Command = subCmd.CommandPath()
	meta.CommandFlags = flags

	meta.CustomProperty = "foo"

	if runErr != nil {
		telemetry.Error(runErr, meta)
	}

	return runErr
}

func getSubcommand( //nolint:nonamedreturns
	cmd *cobra.Command,
	args []string,
) (subCmd *cobra.Command, flags []string, err error) {
	if cmd.TraverseChildren {
		subCmd, _, err = cmd.Traverse(args)
	} else {
		subCmd, _, err = cmd.Find(args)
	}

	if err != nil {
		return nil, nil, err
	}

	subCmd.Flags().Visit(func(f *pflag.Flag) {
		flags = append(flags, "--"+f.Name)
	})

	sort.Strings(flags)

	return subCmd, flags, err
}
