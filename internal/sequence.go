package internal

type Sequence struct {
	current int
}

func NewSequence() *Sequence {
	return &Sequence{
		current: 0,
	}
}

func (s *Sequence) GetNext() int {
	defer func() { s.current++ }()
	return s.current

}
