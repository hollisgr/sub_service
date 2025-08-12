package subscription

import (
	"context"
	"fmt"
	"main/internal/db"
	"main/internal/model"
	"main/pkg/logger"
	"strings"
	"time"
)

type SubscriptionService struct {
	Storage db.Storage
	Logger  *logger.Logger
}

func NewService(s db.Storage, logger *logger.Logger) SubscriptionInterface {
	return &SubscriptionService{
		Storage: s,
		Logger:  logger,
	}
}

func (s *SubscriptionService) Save(ctx context.Context, sub model.Subscription) (id int, err error) {
	id, err = s.Storage.Save(ctx, s.mapperToDTO(sub))
	if err != nil {
		s.Logger.Errorln(err)
		return id, err
	}
	return id, nil
}

func (s *SubscriptionService) Load(ctx context.Context, subID int) (sub model.Subscription, err error) {
	dto, err := s.Storage.Load(ctx, subID)
	if err != nil {
		s.Logger.Errorln(err)
		return sub, err
	}
	return s.mapperToSub(dto), nil
}

func (s *SubscriptionService) LoadList(ctx context.Context) (subs []model.Subscription, err error) {
	dtos, err := s.Storage.LoadList(ctx)
	if err != nil {
		s.Logger.Errorln(err)
		return subs, err
	}
	for _, dto := range dtos {
		sub := s.mapperToSub(dto)
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, subID int) (err error) {
	err = s.Storage.Delete(ctx, subID)
	if err != nil {
		s.Logger.Errorln(err)
		return err
	}
	return nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub model.Subscription) (err error) {
	err = s.Storage.Update(ctx, s.mapperToDTO(sub))
	if err != nil {
		s.Logger.Errorln(err)
		return err
	}
	return nil
}

func (s *SubscriptionService) Cost(ctx context.Context, data model.CostRequest) (cost int, err error) {
	dto := s.mapperCostToDTO(data)
	if dto.StartDate.IsZero() || dto.EndDate.IsZero() {
		return cost, fmt.Errorf("invalid date, MM-YYYY format required")
	}
	cost, err = s.Storage.Cost(ctx, dto)
	if err != nil {
		s.Logger.Errorln(err)
		return cost, err
	}
	return cost, nil
}

func (s *SubscriptionService) convertStringToDate(str string) (date time.Time) {
	minDate := "01-1999"
	if len(str) == 0 || str <= minDate {
		return date
	}
	split := strings.Split(str, "-")
	dateStr := fmt.Sprintf("%s-%s-01", split[1], split[0])
	date, _ = time.Parse("2006-01-02", dateStr)

	return date
}

func (s *SubscriptionService) convertDateToString(date time.Time) (str string) {
	if date.IsZero() {
		return ""
	}
	year, month, _ := date.Date()
	return fmt.Sprintf("%02d-%d", month, year)
}

func (s *SubscriptionService) mapperToDTO(sub model.Subscription) model.SubscriptionDTO {
	return model.SubscriptionDTO{
		Id:          sub.Id,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserId:      sub.UserId,
		StartDate:   s.convertStringToDate(sub.StartDate),
		EndDate:     s.convertStringToDate(sub.EndDate),
	}
}

func (s *SubscriptionService) mapperToSub(dto model.SubscriptionDTO) model.Subscription {
	return model.Subscription{
		Id:          dto.Id,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      dto.UserId,
		StartDate:   s.convertDateToString(dto.StartDate),
		EndDate:     s.convertDateToString(dto.EndDate),
	}
}

func (s *SubscriptionService) mapperCostToDTO(data model.CostRequest) model.CostDTO {
	return model.CostDTO{
		UserId:      data.UserId,
		ServiceName: data.ServiceName,
		StartDate:   s.convertStringToDate(data.StartDate),
		EndDate:     s.convertStringToDate(data.EndDate),
	}
}
