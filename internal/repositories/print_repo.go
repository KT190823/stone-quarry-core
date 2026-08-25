package repositories

// PrintTemplateRepo persists printable document templates.
type PrintTemplateRepo struct {
	*BaseRepo
}

func NewPrintTemplateRepo() *PrintTemplateRepo {
	return &PrintTemplateRepo{BaseRepo: NewBaseRepo("print_templates", "id")}
}
