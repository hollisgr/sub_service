package db

import (
	"context"
	"fmt"
	"main/internal/model"
	"main/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	conn   *pgxpool.Pool
	logger *logger.Logger
}

func NewDataBase(pool *pgxpool.Pool, logger *logger.Logger) Storage {
	return &db{
		conn:   pool,
		logger: logger,
	}
}

func (d *db) Save(ctx context.Context, dto model.SubscriptionDTO) (id int, err error) {
	query := `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id
	`
	err = d.conn.QueryRow(ctx, query, dto.ServiceName, dto.Price, dto.UserId,
		dto.StartDate, dto.EndDate).Scan(&id)
	if id == 0 || err != nil {
		return id, fmt.Errorf("database error, failed to save sub: %v", err)
	}
	return id, nil
}

func (d *db) Load(ctx context.Context, subID int) (dto model.SubscriptionDTO, err error) {
	query := `
		SELECT 
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM 
			subscriptions
		WHERE
			id = $1
	`
	row := d.conn.QueryRow(ctx, query, subID)
	err = row.Scan(&dto.Id, &dto.ServiceName, &dto.Price, &dto.UserId,
		&dto.StartDate, &dto.EndDate)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (d *db) LoadList(ctx context.Context) (dtoList []model.SubscriptionDTO, err error) {
	query := `
		SELECT 
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM 
			subscriptions
	`
	rows, err := d.conn.Query(ctx, query)
	if err != nil {
		return dtoList, err
	}
	for rows.Next() {
		dto := model.SubscriptionDTO{}
		rows.Scan(&dto.Id, &dto.ServiceName, &dto.Price, &dto.UserId,
			&dto.StartDate, &dto.EndDate)
		dtoList = append(dtoList, dto)
	}

	if len(dtoList) == 0 {
		return dtoList, fmt.Errorf("database error, sub list is empty")
	}
	return dtoList, nil
}

func (d *db) Delete(ctx context.Context, subID int) (err error) {
	query := `
		DELETE 
		FROM 
			subscriptions
		WHERE
			id = $1
		RETURNING 
			id
	`
	var tempId int
	d.conn.QueryRow(ctx, query, subID).Scan(&tempId)
	if tempId != subID {
		return fmt.Errorf("database delete error, sub not found")
	}
	return nil
}

func (d *db) Update(ctx context.Context, dto model.SubscriptionDTO) (err error) {
	query := `
		UPDATE
			subscriptions
		SET
			service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6
		WHERE
			id = $1
	`
	res, err := d.conn.Exec(ctx, query, dto.Id, dto.ServiceName, dto.Price, dto.UserId,
		dto.StartDate, dto.EndDate)
	if err != nil {
		return fmt.Errorf("database error, failed to update sub: %v", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("database error, no rows updated")
	}
	return nil
}

func (d *db) Cost(ctx context.Context, data model.CostDTO) (cost int, err error) {
	query := `
		SELECT
			sum(price)
		FROM
			subscriptions
		WHERE
			user_id = $1
			AND
			service_name = $2
			AND
			start_date BETWEEN $3 AND $4
	`
	row := d.conn.QueryRow(ctx, query, data.UserId, data.ServiceName, data.StartDate, data.EndDate)
	err = row.Scan(&cost)
	if err != nil {
		return cost, err
	}
	if cost == 0 {
		return cost, fmt.Errorf("database error, failed to cost request")
	}
	return cost, nil
}
