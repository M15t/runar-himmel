package ulidutil

import (
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

// NewService returns ULID service with some handy functions
func NewService() Service {
	s := service{}
	s.entropy = SecureEntropy()

	// Faster but not secure
	// s.entropy = ulid.DefaultEntropy()

	return &s
}

// An ULID is a 16 byte Universally Unique Lexicographically Sortable Identifier
type ULID = ulid.ULID

// Service represents the interface for the ULID utils
type Service interface {
	// NewString returns an ULID as a string, with the current time. Panic on failure.
	NewString() string
	// NewStringWithTime returns an ULID as a string, with the given time. Panic on failure.
	NewStringWithTime(t time.Time) string
	// New returns an ULID with the current time
	New() (ULID, error)
	// MustNew returns an ULID with the current time. Panic on failure.
	MustNew() ULID
	// NewWithTime returns an ULID with given time
	NewWithTime(t time.Time) (ULID, error)
	// MustNewWithTime returns an ULID with given time. Panic on failure.
	MustNewWithTime(t time.Time) ULID
}

// service implements the Service interface
type service struct {
	entropy io.Reader
}

// NewString returns an ULID as a string, with the current time. Panic on failure.
func (s *service) NewString() string {
	return s.MustNew().String()
}

// NewStringWithTime returns an ULID as a string, with the given time. Panic on failure.
func (s *service) NewStringWithTime(t time.Time) string {
	return s.MustNewWithTime(t).String()
}

// New returns a new ULID with the current time
func (s *service) New() (ULID, error) {
	return s.NewWithTime(time.Now())
}

// MustNew returns a new ULID with the current time. Panic on failure.
func (s *service) MustNew() ULID {
	return s.MustNewWithTime(time.Now())
}

// NewWithTime returns a new ULID with given time
func (s *service) NewWithTime(t time.Time) (ULID, error) {
	return ulid.New(ulid.Timestamp(t), s.entropy)
}

// MustNewWithTime returns a new ULID with given time. Panic on failure.
func (s *service) MustNewWithTime(t time.Time) ULID {
	return ulid.MustNew(ulid.Timestamp(t), s.entropy)
}
