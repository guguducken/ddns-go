package errno

import (
	"errors"

	"github.com/guguducken/ddns-go/pkg/utils/poolutils"
)

type Error struct {
	// code is used internally to quickly locate problems
	code string
	// message is the information that will be displayed to the user
	message string

	additionalInfo additionalInfo
}

func (e Error) Error() string {
	builder := poolutils.GetStringBuilder()
	defer poolutils.PutStringBuilder(builder)

	builder.WriteString("code: ")
	builder.WriteString(e.code)
	builder.WriteString(", message: ")
	builder.WriteString(e.message)
	return builder.String()
}

func NewError(code, message string) error {
	return Error{
		code:    code,
		message: message,
	}
}

type additionalInfo map[string]string

func (e Error) Is(in error) bool {
	var ne Error
	if errors.As(in, &ne) {
		return ne.code == e.code
	}
	return false
}

func GetAdditionalInfo(err error) additionalInfo {
	var ne Error
	if errors.As(err, &ne) {
		return ne.additionalInfo
	}
	return nil
}

func OverrideError(err error, options ...Option) error {
	if err == nil {
		return nil
	}
	var ne Error
	if !errors.As(err, &ne) {
		ne = Error{
			code:    "ErrUnknown",
			message: err.Error(),
		}
	}
	for _, opt := range options {
		opt(&ne)
	}
	return ne
}
