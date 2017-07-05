package argument

import (
	"bytes"
	"errors"
	"io"
)

type Scanner struct {
	cmdLine         string
	pos             int
	buf             *bytes.Buffer
	previousStateFn stateFn
	quote           byte
	err             error
}

type stateFn func(s *Scanner) stateFn

func NewScanner(cmdLine string) *Scanner {
	return &Scanner{
		cmdLine: cmdLine,
		buf:     bytes.NewBuffer(make([]byte, 0, 20)),
	}
}

func (s *Scanner) Rest() string {
	return s.cmdLine
}

func (s *Scanner) Token() string {
	return s.buf.String()
}

func (s *Scanner) Scan() (err error) {
	if err != nil {
		return s.err
	} else if s.cmdLine == "" {
		return ErrUnexpectedEnd
	}
	s.reset()
	s.run()
	s.cmdLine = s.cmdLine[s.pos:]
	if s.err != nil {
		return s.err
	}
	return nil
}

func (s *Scanner) reset() {
	s.buf.Reset()
	s.pos = 0
}

func (s *Scanner) run() {
	for state := scanText; state != nil; {
		state = state(s)
	}
}

func (s *Scanner) get() (byte, error) {
	if s.pos >= len(s.cmdLine) {
		return 0, io.EOF
	}
	return s.cmdLine[s.pos], nil
}

func (s *Scanner) writePrev(n int) {
	s.buf.WriteString(s.cmdLine[s.pos-n: s.pos])
}

func (s *Scanner) writeFrom(ppos int) {
	s.buf.WriteString(s.cmdLine[ppos:s.pos])
}

func scanText(s *Scanner) stateFn {
	ppos := s.pos
	for {
		c, err := s.get()
		if err != nil {
			s.writeFrom(ppos)
			return nil
		}
		switch c {
		case '\\':
			s.writeFrom(ppos)
			s.previousStateFn = scanText
			s.pos++
			return scanBackslash
		case '"', '\'':
			s.writeFrom(ppos)
			s.pos++
			s.quote = c
			return scanQuote
		case ' ':
			s.writeFrom(ppos)
			s.pos++
			return nil
		}
		s.pos++
	}
}

var (
	ErrUnexpectedBackspace = errors.New("You need to have a character after the backspace.")
	ErrUnexpectedEnd = errors.New("I need moar.")
)

// TODO can add special escape sequences, idk if I'll ever need to. On the fence about these semantics. If I choose not to and remove s.err, then I'll change Scanner.Scan() to return a boolean instead of an error
func scanBackslash(s *Scanner) stateFn {
	c, err := s.get()
	if err != nil {
		s.err = ErrUnexpectedBackspace
		return nil
	}
	s.buf.WriteByte(c)
	s.pos++
	return s.previousStateFn
}

func scanQuote(s *Scanner) stateFn {
	ppos := s.pos
	for {
		c, err := s.get()
		if err != nil {
			s.writeFrom(ppos)
			return nil
		}
		switch c {
		case '\\':
			s.writeFrom(ppos)
			s.previousStateFn = scanQuote
			s.pos++
			return scanBackslash
		case s.quote:
			s.writeFrom(ppos)
			s.pos++
			return scanText
		}
		s.pos++
	}
}