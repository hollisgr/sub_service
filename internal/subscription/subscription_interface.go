package subscription

import (
	"context"
	"main/internal/model"
)

type SubscriptionInterface interface {
	Save(ctx context.Context, sub model.Subscription) (id int, err error)
	Load(ctx context.Context, subID int) (sub model.Subscription, err error)
	LoadList(ctx context.Context) (subs []model.Subscription, err error)
	Delete(ctx context.Context, subID int) (err error)
	Update(ctx context.Context, sub model.Subscription) (err error)
	Cost(ctx context.Context, data model.CostRequest) (cost int, err error)
}
