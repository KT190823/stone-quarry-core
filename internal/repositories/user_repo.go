package repositories

type UserRepo struct {
	*BaseRepo
}

func NewUserRepo() *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo("users", "id")}
}

type UserRoleRepo struct {
	*BaseRepo
}

func NewUserRoleRepo() *UserRoleRepo {
	return &UserRoleRepo{BaseRepo: NewBaseRepo("user_roles", "id")}
}

type UserLogRepo struct {
	*BaseRepo
}

func NewUserLogRepo() *UserLogRepo {
	return &UserLogRepo{BaseRepo: NewBaseRepo("user_logs", "id")}
}

type ReportRepo struct {
	*BaseRepo
}

func NewReportRepo() *ReportRepo {
	return &ReportRepo{BaseRepo: NewBaseRepo("reports", "id")}
}

type SettingRepo struct {
	*BaseRepo
}

func NewSettingRepo() *SettingRepo {
	return &SettingRepo{BaseRepo: NewBaseRepo("settings", "id")}
}
