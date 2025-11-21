package dto

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

const (
	TEAM_EXISTS  = "TEAM_EXISTS"
	PR_EXISTS    = "PR_EXISTS"
	PR_MERGED    = "PR_MERGED"
	NOT_ASSIGNED = "NOT_ASSIGNED"
	NO_CANDIDATE = "NO_CANDIDATE"
	NOT_FOUND    = "NOT_FOUND"
)

func ErrTeamExists(message ...string) *ErrorResponse {
	msg := "Team already exists"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    TEAM_EXISTS,
			Message: msg,
		},
	}
}

func ErrPRExists(message ...string) *ErrorResponse {
	msg := "Pull request already exists"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    PR_EXISTS,
			Message: msg,
		},
	}
}

func ErrPRMerged(message ...string) *ErrorResponse {
	msg := "Pull request already merged"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    PR_MERGED,
			Message: msg,
		},
	}
}

func ErrNotAssigned(message ...string) *ErrorResponse {
	msg := "Reviewer is not assigned to this PR"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    NOT_ASSIGNED,
			Message: msg,
		},
	}
}

func ErrNoCandidate(message ...string) *ErrorResponse {
	msg := "No active candidate to review PR in team"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    NO_CANDIDATE,
			Message: msg,
		},
	}
}

func ErrTeamNotFound(message ...string) *ErrorResponse {
	msg := "Team not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    NOT_FOUND,
			Message: msg,
		},
	}
}

func ErrUserNotFound(message ...string) *ErrorResponse {
	msg := "User not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    NOT_FOUND,
			Message: msg,
		},
	}
}

func ErrPRNotFound(message ...string) *ErrorResponse {
	msg := "Pull request not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    NOT_FOUND,
			Message: msg,
		},
	}
}
