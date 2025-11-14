package v1

import (
	"avito-internship/pkg/e"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorResponseDetail `json:"error"`
}

type ErrorResponseDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorResponseDetail{
			Code:    code,
			Message: message,
		},
	}
}

func ToHTTPResponse(err error) (int, string, string) {
	switch {
	case errors.Is(err, e.ErrUserNotFound),
		errors.Is(err, e.ErrTeamNotFound),
		errors.Is(err, e.ErrUnauthorized),
		errors.Is(err, e.ErrStatusNotFound),
		errors.Is(err, e.ErrPRNotFound):
		return http.StatusNotFound, e.NOT_FOUND, e.ErrResourceNotFound.Error()
	case errors.Is(err, e.ErrTeamIsExists):
		return http.StatusBadRequest, e.TEAM_EXISTS, e.ErrTeamIsExists.Error()
	case errors.Is(err, e.ErrPRIsExists):
		return http.StatusBadRequest, e.PR_EXISTS, e.ErrPRIsExists.Error()
	case errors.Is(err, e.ErrPrMerged):
		return http.StatusConflict, e.PR_MERGED, e.ErrPrMerged.Error()
	case errors.Is(err, e.ErrPrReviewerNotAssigned):
		return http.StatusConflict, e.NOT_ASSIGNED, e.ErrPrReviewerNotAssigned.Error()
	case errors.Is(err, e.ErrPrNoCandidate):
		return http.StatusConflict, e.NO_CANDIDATE, e.ErrPrNoCandidate.Error()
	case errors.Is(err, e.ErrEmptyMembers):
		return http.StatusBadRequest, e.BAD_REQUEST, e.ErrEmptyMembers.Error()
	default:
		return http.StatusInternalServerError, e.SERVER_ERR, e.ErrInternalServerError.Error()
	}
}
