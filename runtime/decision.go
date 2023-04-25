package runtime

import "fmt"

type Decision interface {
	Allow() bool
}

var (
	_ Decision = &DecisionError{}
	_ Decision = &DecisionReject{}
	_ Decision = &DecisionPass{}
)

type Trace struct {
	Message string
	Context map[string]interface{}
}

type Undecided struct {
	traces []Trace
}

// Append information to the decisions trace to help explain the action.
func (d *Undecided) Trace(msg string) {
	d.traces = append(d.traces, Trace{Message: msg})
}

// Trace with fmt-formatted message.
func (d *Undecided) Tracef(format string, args ...interface{}) {
	d.Trace(fmt.Sprintf(format, args...))
}

// Append information to the decision trace with extra key:value fields.
func (d *Undecided) TraceWithContext(msg string, ctx map[string]interface{}) {
	d.traces = append(d.traces, Trace{Message: msg, Context: ctx})
}

// Mark the decision as passed, locking it.
func (d *Undecided) Pass() DecisionPass {
	d.Trace("changed state to decision pass")

	return DecisionPass{traces: d.traces}
}

// Mark the decision as rejected, locking it.
func (d *Undecided) Reject() DecisionReject {
	d.Trace("changed state to decision reject")

	return DecisionReject{traces: d.traces}
}

// Mark the decision as error, locking it.
func (d *Undecided) Error(err error) DecisionError {
	d.Trace(fmt.Sprintf("an error occurred: %v", err))

	return DecisionError{
		traces: d.traces,
		err:    err,
	}
}

// Decision resulting in an error result.
type DecisionError struct {
	err    error
	traces []Trace
}

func (d *DecisionError) Allow() bool {
	return false
}

// Decision resulting in a rejection result.
type DecisionReject struct {
	traces []Trace
}

func (d *DecisionReject) Allow() bool {
	return false
}

// Decision resulting in a pass result.
type DecisionPass struct {
	traces []Trace
}

func (d *DecisionPass) Allow() bool {
	return true
}
