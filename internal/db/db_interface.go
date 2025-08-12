package db

import (
	"context"
	"main/internal/model"
)

type Storage interface {
	Save(ctx context.Context, sub model.SubscriptionDTO) (id int, err error)
	Load(ctx context.Context, subID int) (sub model.SubscriptionDTO, err error)
	LoadList(ctx context.Context) (subList []model.SubscriptionDTO, err error)
	Delete(ctx context.Context, subID int) (err error)
	Update(ctx context.Context, sub model.SubscriptionDTO) (err error)
	Cost(ctx context.Context, data model.CostDTO) (cost int, err error)
}
