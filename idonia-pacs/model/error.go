package model

import "encoding/json"

type Error struct {
	Code              uint32 `json:"code,omitempty"`
	Message           string `json:"message,omitempty"`
	AdditionalMessage string `json:"additionalMessage,omitempty"`
}

func (error Error) String() string {
	errorString, _ := json.Marshal(error)
	return string(errorString)
}

func (error Error) Error() string {
	return error.String()
}

var (
	ErrUndefined = Error{Code: 0, Message: "unexpected error"}

	ErrUser                     = Error{Code: 100, Message: "user undefined error"}
	ErrUserInvalidCredentials   = Error{Code: 101, Message: "invalid username or password"}
	ErrUserInvalidBearer        = Error{Code: 102, Message: "invalid or missing bearer token"}
	ErrUserDestinationsNotFound = Error{Code: 103, Message: "user does not have destinations"}
	ErrUserMissingData          = Error{Code: 104, Message: "missing data in request"}
	ErrUserInvalidData          = Error{Code: 105, Message: "invalid data in request"}

	ErrPACS           = Error{Code: 200, Message: "PACS undefined error"}
	ErrPACSCFindError = Error{Code: 201, Message: "PACS error on c-find"}

	ErrIdonia = Error{Code: 300, Message: "idonia undefined error"}

	ErrConf                 = Error{Code: 400, Message: "configuration undefined error"}
	ErrConfStepTypeMistyped = Error{Code: 401, Message: "step name mistyped"}

	ErrServer   = Error{Code: 500, Message: "server undefined error"}
	ErrServerDB = Error{Code: 501, Message: "server database error"}
)
