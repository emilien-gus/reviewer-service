package models

const (
	CodeTeamExists          = "TEAM_EXISTS"
	CodePRExists            = "PR_EXISTS"
	CodePRMerged            = "PR_MERGED"
	CodeNotAssigned         = "NOT_ASSIGNED"
	CodeNoCandidate         = "NO_CANDIDATE"
	CodeNotFound            = "NOT_FOUND"
	CodeInvalidRequest      = "INVALID_REQUEST"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	var resp ErrorResponse
	resp.Error.Code = code
	resp.Error.Message = message
	return resp
}
