package loglet

type SequenceNumberer interface {
	Max() SequenceNumberer
	Invalid() SequenceNumberer
	Oldest() SequenceNumberer
	Next() SequenceNumberer
	Prev() SequenceNumberer
}
