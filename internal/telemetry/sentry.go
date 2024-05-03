package telemetry

import (
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"

	"github.com/intility/minctl/internal/build"
)

const (
	transportTimeout = time.Second * 2
)

var ExecutionID = newEventID()

func initSentryClient(appName string) bool {
	if appName == "" {
		panic("telemetry.Start: app name is empty")
	}

	// This will be set by the build system.
	//noinspection GoBoolExpressions
	if build.SentryDSN == "" {
		return false
	}

	transport := sentry.NewHTTPTransport()
	transport.Timeout = transportTimeout
	environment := "production"

	if build.IsDev {
		environment = "development"
	}

	err := sentry.Init(sentry.ClientOptions{ //nolint:exhaustruct
		Dsn:              build.SentryDSN,
		Environment:      environment,
		Release:          appName + "@" + build.Version,
		Transport:        transport,
		TracesSampleRate: 1,
		BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
			// redact the hostname, which the SDK automatically adds
			event.ServerName = ""
			return event
		},
	})

	return err == nil
}

func newSentryException(errToLog error) []sentry.Exception { //nolint:funlen,cyclop
	errMsg := errToLog.Error()
	binPkg := ""
	modPath := ""

	if bld, ok := debug.ReadBuildInfo(); ok {
		binPkg = bld.Path
		modPath = bld.Main.Path
	}

	// Unwrap in a loop to get the most recent stack trace. stFunc is set to a
	// function that can generate a stack trace for the most recent error. This
	// avoids computing the full stack trace for every error.
	var stFunc func() []runtime.Frame

	errType := "Generic Error"

	for {
		if t := exportedErrType(errToLog); t != "" {
			errType = t
		}

		//nolint:errorlint
		switch stackErr := errToLog.(type) {
		// If the error implements the StackTrace method in the redact package, then
		// prefer that. The Sentry SDK gets some things wrong when guessing how
		// to extract the stack trace.
		case interface{ StackTrace() []runtime.Frame }:
			stFunc = stackErr.StackTrace
		// Otherwise use the pkg/errors StackTracer interface.
		case interface{ StackTrace() errors.StackTrace }:
			// Normalize the pkgs/errors.StackTrace type to a slice of runtime.Frame.
			stFunc = func() []runtime.Frame {
				pkgStack := stackErr.StackTrace()
				pc := make([]uintptr, len(pkgStack))

				for i := range pkgStack {
					pc[i] = uintptr(pkgStack[i])
				}

				frameIter := runtime.CallersFrames(pc)
				frames := make([]runtime.Frame, 0, len(pc))

				for {
					frame, more := frameIter.Next()
					frames = append(frames, frame)

					if !more {
						break
					}
				}

				return frames
			}
		}

		uw := errors.Unwrap(errToLog)
		if uw == nil {
			break
		}

		errToLog = uw
	}

	ex := []sentry.Exception{{Type: errType, Value: errMsg}} //nolint:exhaustruct
	if stFunc != nil {
		ex[0].Stacktrace = newSentryStack(stFunc(), binPkg, modPath)
	}

	return ex
}

func newSentryStack(frames []runtime.Frame, binPkg, modPath string) *sentry.Stacktrace {
	stack := &sentry.Stacktrace{ //nolint:exhaustruct
		Frames: make([]sentry.Frame, len(frames)),
	}

	for i, frame := range frames {
		pkgName, funcName := splitPkgFunc(frame.Function)

		// The entrypoint has the full function name "main.main". Replace the
		// package name with its full package path to make it easier to find.
		if pkgName == "main" {
			pkgName = binPkg
		}

		// The file path will be absolute unless the binary was built with -trimpath
		// (which releases should be). Absolute paths make it more difficult for
		// Sentry to correctly group errors, but there's no way to infer a relative
		// path from an absolute path at runtime.
		var absPath, relPath string
		if filepath.IsAbs(frame.File) {
			absPath = frame.File
		} else {
			relPath = frame.File
		}

		// Reverse the frames - Sentry wants the most recent call first.
		stack.Frames[len(frames)-i-1] = sentry.Frame{ //nolint:exhaustruct
			Function: funcName,
			Module:   pkgName,
			Filename: relPath,
			AbsPath:  absPath,
			Lineno:   frame.Line,
			InApp:    strings.HasPrefix(frame.Function, modPath) || pkgName == binPkg,
		}
	}

	return stack
}

// exportedErrType returns the underlying type name of err if it's exported.
// Otherwise, it returns an empty string.
func exportedErrType(err error) string {
	t := reflect.TypeOf(err)
	if t == nil {
		return ""
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	name := t.Name()

	if r, _ := utf8.DecodeRuneInString(name); unicode.IsUpper(r) {
		return t.String()
	}

	return ""
}

// splitPkgFunc splits a fully-qualified function or method name into its
// package path and base name components.
func splitPkgFunc(name string) (string, string) {
	// Using the following fully-qualified function name as an example:
	// github.com/intility/icpctl/pkg/wizard.(*Wizard).Run
	//
	// dir = github.com/intility/icpctl/pkg/
	// base = wizard.(*Wizard).Run
	dir, base := path.Split(name)

	// pkgName = wizard
	// fn = (*Wizard).Run
	pkgName, fn, _ := strings.Cut(base, ".")

	// pkgPath = github.com/intility/icpctl/pkg/wizard
	// funcName = (*Wizard).Run
	return dir + pkgName, fn
}

// bufferSentryEvent buffers a Sentry event to disk so that Report can upload it
// later.
func bufferSentryEvent(event *sentry.Event) {
	bufferEvent(filepath.Join(sentryBufferDir, string(event.EventID)+".json"), event)
}