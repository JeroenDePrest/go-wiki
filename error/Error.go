package error

type Error struct {
	Message string `json:"message"`
}

func Create(message string) Error{
	return Error{message}
}
