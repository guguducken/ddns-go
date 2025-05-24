package errno

import (
	"errors"
	"fmt"
	"testing"
)

func Test_ErrorIs(t *testing.T) {
	e1 := Error{code: "e1", message: "e1", additionalInfo: additionalInfo{"aaaaaa": "bbbbbb"}}
	e2 := Error{code: "e1", message: "e2"}
	fmt.Printf("errors.Is(e1, e2): %v\n", errors.Is(e1, e2))
	fmt.Printf("e1.Error(): %v\n", e1.Error())
}
