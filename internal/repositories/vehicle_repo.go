package repositories

type VehicleRepo struct {
	*BaseRepo
}

func NewVehicleRepo() *VehicleRepo {
	return &VehicleRepo{BaseRepo: NewBaseRepo("vehicles", "id")}
}
