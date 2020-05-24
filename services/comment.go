package service

type Comment struct {
	UserId   int64  `json:"user_id"`
	Content  string `json:"content"`
	Children map[string]Comment
}
