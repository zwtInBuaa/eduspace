package models

// casbin 增删改角色权限所需参数
type Permission struct {
	UserID uint   `json:"user_id"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

// casbin 增删改角色权限所需参数
type RolePermission struct {
	Role   string `json:"role"`
	UserID uint   `json:"user_id"`
}
