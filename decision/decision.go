package decision

// Decision represents a decision made by an authorizer.
type Decision struct {
	Allow  bool
	Reason string
}

// Skip is equivalent to saying "no", or "deny". It doesn't affect other rules.
func Skip(reason string) Decision {
	return Decision{
		Allow:  false,
		Reason: reason,
	}
}

// Allow is equivalent to saying "yes". When returned, no more rules will be evaluated.
func Allow(reason string) Decision {
	return Decision{
		Allow:  true,
		Reason: reason,
	}
}

// Error is equivalent to saying "Skip" except it uses `err.Error()` as the reason.
func Error(err error) Decision {
	return Decision{
		Allow:  false,
		Reason: err.Error(),
	}
}
