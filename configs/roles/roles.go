package roles

const (
	RoleAdmin  = "ADMIN"
	RoleNormal = "NORMAL"
	RoleBasic  = "BASIC"
)

func GetRolesAvailable() []string {
	return []string{RoleAdmin, RoleNormal, RoleBasic}
}
