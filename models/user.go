package models

// Struct đại diện cho người dùng
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Sử dụng mật khẩu hash trong thực tế
}
