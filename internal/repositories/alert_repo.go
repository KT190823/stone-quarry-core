package repositories

type AlertRepo struct {
	*BaseRepo
}

func NewAlertRepo() *AlertRepo {
	return &AlertRepo{BaseRepo: NewBaseRepo("alerts", "id")}
}
