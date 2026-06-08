package dto

type UserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin manager user"`
}
