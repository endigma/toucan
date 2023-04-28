package decision

// Decision represents a decision made by an authorizer.
type Decision struct {
	Allow  bool
	Reason string
}

// False is equivalent to saying "no", or "deny". It doesn't affect other rules.
func False(reason string) Decision {
	return Decision{
		Allow:  false,
		Reason: reason,
	}
}

// True is equivalent to saying "yes". When returned, no more rules will be evaluated.
func True(reason string) Decision {
	return Decision{
		Allow:  true,
		Reason: reason,
	}
}

// Error is equivalent to saying "False" except it uses `err.Error()` as the reason.
func Error(err error) Decision {
	return Decision{
		Allow:  false,
		Reason: err.Error(),
	}
}
