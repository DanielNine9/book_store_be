package dtos

type CategoryRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}