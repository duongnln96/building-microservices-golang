package helper

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  error       `json:"error"`
	Data    interface{} `json:"data"`
}

func BuildResponse(status bool, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func BuildErrorResponse(message string, err error, data interface{}) Response {
	return Response{
		Status:  false,
		Message: message,
		Errors:  err,
		Data:    data,
	}
}
