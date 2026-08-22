package repositories

type DebtRepo struct {
	*BaseRepo
}

func NewDebtRepo() *DebtRepo {
	return &DebtRepo{BaseRepo: NewBaseRepo("payments_debt", "id")}
}

type ReconcileRepo struct {
	*BaseRepo
}

func NewReconcileRepo() *ReconcileRepo {
	return &ReconcileRepo{BaseRepo: NewBaseRepo("payments_reconcile", "id")}
}
