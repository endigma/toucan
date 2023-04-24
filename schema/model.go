package schema

import (
	"errors"
	"fmt"
	"regexp"
)

var modelRegex = regexp.MustCompile(`^([A-z][\w|.|/]+)(?:\.)([A-Z][a-zA-Z0-9_-]*)$`)

type Model struct {
	Path string `validate:"required,validQualPath"`
	Name string `validate:"required,validName,validQualName"`
}

func (m *Model) String() string {
	return fmt.Sprintf("%s.%s", m.Path, m.Name)
}

func (m *Model) Tuple() (string, string) {
	return m.Path, m.Name
}

func (m *Model) Validate() error {
	return fmt.Errorf("invalid model: %w", validate.Struct(m))
}

// encoding.TextMarshaller.
func (m *Model) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf(`%s.%s`, m.Path, m.Name)), nil
}

var ErrInvalidModel = errors.New("invalid model string")

// encoding.TextUnmarshaller.
func (m *Model) UnmarshalText(text []byte) error {
	str := string(text)

	matches := modelRegex.FindStringSubmatch(str)

	if len(matches) == 0 {
		return fmt.Errorf("%w: %s", ErrInvalidModel, str)
	}

	m.Path = matches[1]
	m.Name = matches[2]

	return nil
}
