package response

type Success struct {
	Message   any  `json:"message"`
	Data      any  `json:"data"`
	IsSuccess bool `json:"isSuccess"`
}

func Ok(data any) Success {
	return Success{
		Message:   nil,
		Data:      data,
		IsSuccess: true,
	}
}

type Error struct {
	Message   string `json:"message"`
	Data      any    `json:"data"`
	IsSuccess bool   `json:"isSuccess"`
}

func Failure(message string) Error {
	return Error{
		Message:   message,
		Data:      nil,
		IsSuccess: false,
	}
}
