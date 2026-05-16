package server

const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodePayloadTooLarge = "PAYLOAD_TOO_LARGE"
)

var (
	ErrValidation = ErrorContract{
		Code:    ErrCodeValidation,
		Message: "Validation error",
	}

	ErrBadRequest = ErrorContract{
		Code:    ErrCodeBadRequest,
		Message: "Bad request",
	}

	ErrUnauthorized = ErrorContract{
		Code:    ErrCodeUnauthorized,
		Message: "Authorization required",
	}

	ErrForbidden = ErrorContract{
		Code:    ErrCodeForbidden,
		Message: "Forbidden",
	}

	ErrNotFound = ErrorContract{
		Code:    ErrCodeNotFound,
		Message: "Resource not found",
	}

	ErrConflict = ErrorContract{
		Code:    ErrCodeConflict,
		Message: "Conflict",
	}

	ErrInternal = ErrorContract{
		Code:    ErrCodeInternal,
		Message: "Internal server error",
	}

	ErrPayloadTooLarge = ErrorContract{
		Code:    ErrCodePayloadTooLarge,
		Message: "Payload too large",
	}
)
