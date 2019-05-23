package model

// InvalidReasonForDeveloper is InvalidReason message for developer.
type InvalidReasonForDeveloper string

// String returns the message in string.
func (p InvalidReasonForDeveloper) String() string {
	return string(p)
}

// DomainModelNameForDeveloper is Model name for developer.
type DomainModelNameForDeveloper string

func (p DomainModelNameForDeveloper) String() string {
	return string(p)
}

// Model name for developer.
const (
	DomainModelNameUserForDeveloper DomainModelNameForDeveloper = "User"
)

// DomainModelNameForUser is Model name for user.
type DomainModelNameForUser string

// String returns in string.
func (p DomainModelNameForUser) String() string {
	return string(p)
}

// Model name for user.
const (
	DomainModelNameUserForUser DomainModelNameForUser = "ユーザー"
)

// PropertyNameForDeveloper is property name for developer.
type PropertyNameForDeveloper string

// String return as string.
func (p PropertyNameForDeveloper) String() string {
	return string(p)
}

// Property name for developer.
const (
	IDPropertyForDeveloper PropertyNameForDeveloper = "id"
)

// PropertyNameForUser is Property name for user.
type PropertyNameForUser string

// String return as string.
func (p PropertyNameForUser) String() string {
	return string(p)
}

// Property name for user.
const (
	IDPropertyForUser PropertyNameForUser = "ID"
)
