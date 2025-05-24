package errno

type Option func(err *Error)

func OverrideMessage(message string) Option {
	return func(err *Error) {
		if len(message) != 0 {
			err.message = message
		}
	}
}

func AppendAdditionalMessage(key string, message string) Option {
	return func(err *Error) {
		if len(message) != 0 {
			if len(err.additionalInfo) == 0 {
				err.additionalInfo = make(AdditionalInfo, 2)
			}
			err.additionalInfo[key] = message
		}
	}
}
