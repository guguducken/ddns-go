package logutil

type Fields []Field
type Field struct {
	Key   string
	Value string
}

func NewField(key, value string) Field {
	return Field{Key: key, Value: value}
}
