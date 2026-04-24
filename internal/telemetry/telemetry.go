package telemetry

import (
	"strings"
)

type Event struct {
	Template   string
	DryRun     bool
	Force      bool
	ProjectDir string
}

type Sink interface {
	Record(Event)
}

type Telemetry struct {
	SessionEnabled bool
	LocalEnabled   bool
	ProjectEnabled bool

	SessionSink Sink
	LocalSink   Sink
	ProjectSink Sink
}

func ParseModes(input string) (session, local, project bool) {
	input = strings.ToLower(input)
	if input == "" {
		return true, false, false // defaults
	}

	// Normalize separators
	tokens := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == '.'
	})

	for _, t := range tokens {
		switch t {
		case "none", "off":
			return false, false, false
		case "all", "everything":
			return true, true, true
		case "session":
			session = true
		case "local":
			local = true
		case "project":
			project = true
		}
	}

	return
}

func (t Telemetry) Record(ev Event) {
	if t.SessionEnabled && t.SessionSink != nil {
		t.SessionSink.Record(ev)
	}
	if t.LocalEnabled && t.LocalSink != nil {
		t.LocalSink.Record(ev)
	}
	if t.ProjectEnabled && t.ProjectSink != nil {
		t.ProjectSink.Record(ev)
	}
}

// New creates a new Telemetry instance with the specified configuration.
func New(enabled, session, local, project bool) Telemetry {
	return Telemetry{
		SessionEnabled: enabled && session,
		LocalEnabled:   enabled && local,
		ProjectEnabled: enabled && project,
		SessionSink:    &SessionSink{},
		LocalSink:      &LocalSink{},
		ProjectSink:    &ProjectSink{},
	}
}
