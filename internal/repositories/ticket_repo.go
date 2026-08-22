package repositories

type TicketRepo struct {
	*BaseRepo
}

func NewTicketRepo() *TicketRepo {
	return &TicketRepo{BaseRepo: NewBaseRepo("tickets", "id")}
}
