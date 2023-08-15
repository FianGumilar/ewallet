package template

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repository struct {
	db *sql.DB
}

func NewTemplateRepository(con *sql.DB) domain.TemplateRepository {
	return &repository{db: con}
}

// FindByCode implements domain.TemplateRepository.
func (r repository) FindByCode(ctx context.Context, code string) (template domain.Template, err error) {
	query := `SELECT * FROM templates WHERE code = ?`

	row := r.db.QueryRowContext(ctx, query, code)
	err = row.Scan(&template.Code, &template.Title, &template.Body)
	if err != nil {
		return template, nil
	}
	return
}
