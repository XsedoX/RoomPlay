package response

type Success struct {
	Message   *string `json:"message"  example:"null" extensions:"x-nullable"`
	Data      any     `json:"data" swaggertype:"object"`
	IsSuccess bool    `json:"isSuccess"`
}

func Ok(data any) Success {
	return Success{
		Message:   nil,
		Data:      data,
		IsSuccess: true,
	}
}

type Error struct {
	Message   string  `json:"message"`
	Data      *string `json:"data" example:"null" extensions:"x-nullable"`
	IsSuccess bool    `json:"isSuccess" example:"false"`
}

func Failure(message string) Error {
	return Error{
		Message:   message,
		Data:      nil,
		IsSuccess: false,
	}
}
