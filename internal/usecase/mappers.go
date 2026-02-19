package usecase

import "github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"

// toUserResponse converts a domain User entity (with preloaded Role) into
// a UserResponse DTO suitable for API output.
func toUserResponse(u *domain.User) *domain.UserResponse {
	role := ""
	if u.Role != nil {
		role = u.Role.Code
	}
	return &domain.UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FullName:        u.FullName,
		BirthDate:       u.BirthDate,
		Phone:           u.Phone,
		Status:          u.Status,
		EmailVerifiedAt: u.EmailVerifiedAt,
		Role:            role,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
