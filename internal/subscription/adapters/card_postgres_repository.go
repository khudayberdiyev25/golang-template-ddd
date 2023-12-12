package adapters

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/common/postgres"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer_card"
	"golang.org/x/net/context"
	"time"
)

type dbCard struct {
	guid       uuid.UUID
	customerID uuid.UUID
	cardNumber string
	expiry     string
	cardType   string
	token      string
	status     string
	maskedPan  string
	owner      string
	own        bool
	provider   string
	createdAt  time.Time
	updatedAt  time.Time
}

type cardPostgresRepo struct {
	db        *postgres.PostgresDB
	factory   customer_card.Factory
	tableName string
}

func NewDBCardRepository(factory customer_card.Factory, db *postgres.PostgresDB) customer_card.Repository {
	return &cardPostgresRepo{
		db:        db,
		factory:   factory,
		tableName: "customer_card",
	}
}

func (r cardPostgresRepo) FindOneByGuid(ctx context.Context, guid uuid.UUID) (*customer_card.Card, error) {
	query := r.buildCardSelectQuery()

	query.Where(sq.Eq{"guid": guid})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var c dbCard
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&c.guid,
		&c.customerID,
		&c.cardNumber,
		&c.expiry,
		&c.cardType,
		&c.token,
		&c.status,
		&c.maskedPan,
		&c.owner,
		&c.own,
		&c.provider,
		&c.createdAt,
		&c.updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.marshalCardToDomain(&c)
}

func (r cardPostgresRepo) FindOne(ctx context.Context, model *customer_card.ReadModel) (*customer_card.Card, error) {
	query := r.buildCardSelectQuery()

	query = r.applyQueryFilters(query, model)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var c dbCard
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&c.guid,
		&c.customerID,
		&c.cardNumber,
		&c.expiry,
		&c.cardType,
		&c.token,
		&c.status,
		&c.maskedPan,
		&c.owner,
		&c.own,
		&c.provider,
		&c.createdAt,
		&c.updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return r.marshalCardToDomain(&c)
}

func (r cardPostgresRepo) FindAll(ctx context.Context, model *customer_card.ReadModel) ([]*customer_card.Card, error) {
	var list []*customer_card.Card
	query := r.buildCardSelectQuery()
	query = r.applyQueryFilters(query, model)

	if val, ok := model.Limit(); ok {
		query = query.Limit(val)
	}
	if val, ok := model.Limit(); ok {
		query = query.Offset(val)
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return list, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	defer rows.Close()
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var c dbCard
		err = rows.Scan(
			&c.guid,
			&c.customerID,
			&c.cardNumber,
			&c.expiry,
			&c.cardType,
			&c.token,
			&c.status,
			&c.maskedPan,
			&c.owner,
			&c.own,
			&c.provider,
			&c.createdAt,
			&c.updatedAt,
		)

		if err != nil {
			return list, err
		}

		domain, err := r.marshalCardToDomain(&c)
		if err != nil {
			return list, err
		}

		list = append(list, domain)
	}

	return list, nil
}

func (r cardPostgresRepo) Create(ctx context.Context, card *customer_card.Card) error {
	clauses := map[string]any{
		"guid":        card.Guid(),
		"customer_id": card.CustomerID(),
		"card_number": card.CardNumber(),
		"expiry":      card.Expiry(),
		"card_type":   card.CardType(),
		"status":      card.Status(),
		"masked_pan":  card.MaskedPan(),
		"token":       card.Token(),
		"owner":       card.Owner(),
		"own":         card.Own(),
		"provider":    card.Provider(),
		"created_at":  card.CreatedAt(),
		"updated_at":  card.UpdatedAt(),
	}

	sqlStr, args, err := sq.Insert(r.tableName).SetMap(clauses).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.db.Exec(ctx, sqlStr, args...); err != nil {
		return r.db.Error(err)
	}

	return nil
}

func (r cardPostgresRepo) Update(ctx context.Context, card *customer_card.Card) error {
	clauses := map[string]interface{}{
		"customer_id": card.CustomerID(),
		"card_number": card.CardNumber(),
		"expiry":      card.Expiry(),
		"card_type":   card.CardType(),
		"status":      card.Status(),
		"masked_pan":  card.MaskedPan(),
		"token":       card.Token(),
		"owner":       card.Owner(),
		"own":         card.Own(),
		"provider":    card.Provider(),
		"updated_at":  card.UpdatedAt(),
	}

	sqlStr, args, err := sq.Update(r.tableName).
		SetMap(clauses).
		Where(sq.Eq{"guid": card.Guid()}).
		ToSql()

	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no sql rows affected")
	}

	return nil
}

func (r cardPostgresRepo) buildCardSelectQuery() sq.SelectBuilder {
	query := sq.Select(
		"guid",
		"customer_id",
		"card_number",
		"expiry",
		"card_type",
		"token",
		"status",
		"masked_pan",
		"owner",
		"own",
		"provider",
		"created_at",
		"updated_at",
	).From(r.tableName)

	return query
}

func (r cardPostgresRepo) applyQueryFilters(query sq.SelectBuilder, model *customer_card.ReadModel) sq.SelectBuilder {
	if val, ok := model.CustomerID(); ok {
		query = query.Where(sq.Eq{"customer_id": val})
	}

	if val, ok := model.CardNumber(); ok {
		query = query.Where(sq.Eq{"card_number": val})
	}

	return query
}

func (r cardPostgresRepo) marshalCardToDomain(c *dbCard) (*customer_card.Card, error) {
	return r.factory.UnmarshallCardFromDatabase(
		c.guid,
		c.customerID,
		c.cardNumber,
		c.expiry,
		c.cardType,
		c.token,
		c.status,
		c.maskedPan,
		c.owner,
		c.own,
		c.provider,
		c.createdAt,
		c.updatedAt,
	)
}
