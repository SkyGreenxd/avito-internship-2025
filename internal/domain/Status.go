package domain

import "avito-internship/pkg/e"

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type Status struct {
	Id   int
	Name PRStatus
}

func NewStatus(name PRStatus) Status {
	return Status{
		Name: name,
	}
}

func ParseStatus(s string) (PRStatus, error) {
	switch s {
	case string(OPEN):
		return OPEN, nil
	case string(MERGED):
		return MERGED, nil
	}

	return "", e.ErrInvalidStatus
}
