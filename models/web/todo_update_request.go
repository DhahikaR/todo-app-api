package web

type TodoUpdateRequest struct {
	Id          int    `validate:"required"`
	Title       string `validate:"required,min=2,max=200"`
	Description string `validate:"required"`
	Status      string `validate:"omitempty,oneof=pending done"`
}
