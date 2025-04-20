package recorder

import (
	"context"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/recorder/dnspod"
	"gopkg.in/yaml.v3"
)

type Recorder interface {
	GetType() cons.RecorderType
	ApplyValue(ctx context.Context, ip string) (err error)
	Exit(ctx context.Context) (err error)
}

func NewRecorder(ctx context.Context, t cons.RecorderType, config yaml.Node) (Recorder, error) {
	switch t {
	case cons.RecorderTypeDNSPod:
		return dnspod.NewRecorder(ctx, config)
	}
	return nil, errno.OverrideError(
		errno.ErrInvalidRecorderType,
		errno.AppendAdditionalMessage("type", string(t)),
	)
}
