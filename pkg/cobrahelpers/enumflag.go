package cobrahelpers

import "fmt"

// based on https://github.com/hashicorp/packer/blob/v1.6.5/helper/enumflag/flag.go
// changes:
// - different constructor name
// - default value

type enumFlag struct {
	target  *string
	options []string
}

// New returns a flag.Value implementation for parsing flags with a one-of-a-set value
func NewEnumFlag(target *string, defaultValue string, options ...string) *enumFlag {
	*target = defaultValue

	return &enumFlag{target: target, options: options}
}

func (f *enumFlag) String() string {
	return *f.target
}

func (f *enumFlag) Type() string {
	return "enum"
}

func (f *enumFlag) Set(value string) error {
	for _, v := range f.options {
		if v == value {
			*f.target = value
			return nil
		}
	}

	return fmt.Errorf("expected one of %q", f.options)
}
