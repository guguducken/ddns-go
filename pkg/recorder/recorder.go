package recorder

import (
	"context"
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/guguducken/ddns-go/pkg/cons"
	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/recorder/dnspod"
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

type Recorders []Recorder

func (rs Recorders) ApplyValue(ctx context.Context, ip string) error {
	errs := make([]error, 0, len(rs))
	for _, r := range rs {
		if err := r.ApplyValue(ctx, ip); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (rs Recorders) Exit(ctx context.Context) error {
	errs := make([]error, 0, len(rs))
	for _, r := range rs {
		if err := r.Exit(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
