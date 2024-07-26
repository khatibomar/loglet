package loglet

import (
	"context"
	"math"
)

type (
	TailState struct {
		Offset SequenceNumberer
		State  tailStateEnum
	}

	SendableLogletReadStream[S SequenceNumberer, D Payload] struct {
		Stream chan LogRecord[S, D]
		Done   chan struct{}
	}

	LogletBase[S SequenceNumberer, D Payload] interface {
		CreateReadStream(ctx context.Context, from, to SequenceNumberer) (SendableLogletReadStream[S, D], error)
		Append(ctx context.Context, data []byte) (SequenceNumberer, error)
		LastKnownUnsealedTail() SequenceNumberer
		WatchTail() <-chan TailState
		AppendBatch(ctx context.Context, payloads [][]byte) (SequenceNumberer, error)
		FindTail(ctx context.Context) (TailState, error)
		GetTrimPoint(ctx context.Context) (SequenceNumberer, error)
		Trim(ctx context.Context, trimPoint SequenceNumberer) error
		Seal(ctx context.Context) error
		Read(ctx context.Context, from SequenceNumberer) (LogRecord[S, D], error)
		ReadOpt(ctx context.Context, from SequenceNumberer) (LogRecord[S, D], bool, error)
	}

	Loglet[S SequenceNumberer, D Payload] interface {
		LogletBase[S, D]
	}

	LogletOffset struct {
		Value uint64
	}

	LogRecord[S SequenceNumberer, D Payload] struct {
		Offset S
		Record Record[S, D]
	}

	Record[S SequenceNumberer, D Payload] struct {
		Offset S
		Data   D
	}

	Header struct {
		CreatedAt int64 `json:"created_at"`
	}

	Payload struct {
		Body   []byte
		Header Header
	}
)

var _ SequenceNumberer = LogletOffset{}

type tailStateEnum int

const (
	OpenState tailStateEnum = iota
	SealedState
)

func (l LogletOffset) Invalid() SequenceNumberer {
	return LogletOffset{0}
}

func (l LogletOffset) Max() SequenceNumberer {
	return LogletOffset{math.MaxUint64}
}

func (l LogletOffset) Oldest() SequenceNumberer {
	return LogletOffset{1}
}

func (l LogletOffset) Next() SequenceNumberer {
	return l.Add(1)
}

func (l LogletOffset) Prev() SequenceNumberer {
	panic("unimplemented")
}

func (o LogletOffset) Add(rhs uint64) SequenceNumberer {
	sum := o.Value + rhs
	if sum < rhs {
		sum = math.MaxUint64
	}
	return LogletOffset{sum}
}
