package response

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
	} `json:"error"`
}

func TEAM_EXISTS() *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "TEAM_EXISTS",
			Message: "team_name already exists",
		},
	}
}

func PR_EXISTS() *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "PR_EXISTS",
			Message: "PR id already exists",
		},
	}
}

func PR_MERGED() *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "PR_MERGED",
			Message: "cannot reassign on merged PR",
		},
	}
}

func NOT_ASSIGNED() *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "NOT_ASSIGNED",
			Message: "reviewer is not assigned to this PR",
		},
	}
}

func NO_CANDIDATE() *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "NO_CANDIDATE",
			Message: "no active replacement candidate in team",
		},
	}
}

func NOT_FOUND(message ...string) *ErrorResponse {
	if len(message) > 0 {
		return &ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message,omitempty"`
			}{
				Code:    "NOT_FOUND",
				Message: message[0],
			},
		}
	}
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		},
	}
}

func ERROR(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message,omitempty"`
		}{
			Code:    code,
			Message: message,
		},
	}
}
