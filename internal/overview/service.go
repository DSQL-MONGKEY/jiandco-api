package overview

import (
	"context"
	"sort"
	"time"

	"jianco-api/pkg/apperror"
	"jianco-api/pkg/utils"
)

type Service struct {
	repo     *Repository
	location *time.Location
}

func NewService(repo *Repository, location *time.Location) *Service {
	return &Service{repo: repo, location: location}
}

func (s *Service) Revenue(ctx context.Context, month string) (*Response, error) {
	start, end, err := utils.ParseMonth(month, s.location)
	if err != nil {
		return nil, apperror.BadRequest(err.Error())
	}

	revenue, err := s.repo.Revenue(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	expenses, err := s.repo.Expenses(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	totalOrders, err := s.repo.CountOrdersCreated(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	paidOrders, err := s.repo.CountPaidOrders(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	paidOrderRows, err := s.repo.PaidOrders(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	statusBreakdown, err := s.repo.CountOrdersByStatus(ctx, start, end)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	dailyMap := make(map[string]*DailyRevenue)
	for _, row := range paidOrderRows {
		date := row.PaidAt.In(s.location).Format("2006-01-02")
		entry, ok := dailyMap[date]
		if !ok {
			entry = &DailyRevenue{Date: date}
			dailyMap[date] = entry
		}
		entry.Revenue += row.TotalAmount
		entry.Orders++
	}

	dates := make([]string, 0, len(dailyMap))
	for date := range dailyMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	daily := make([]DailyRevenue, 0, len(dates))
	for _, date := range dates {
		daily = append(daily, *dailyMap[date])
	}

	average := int64(0)
	if paidOrders > 0 {
		average = revenue / paidOrders
	}

	return &Response{
		Month:                start.Format(utils.MonthLayout),
		Revenue:              revenue,
		Expenses:             expenses,
		NetProfit:            revenue - expenses,
		TotalOrders:          totalOrders,
		PaidOrders:           paidOrders,
		AverageOrderValue:    average,
		OrderStatusBreakdown: statusBreakdown,
		Daily:                daily,
	}, nil
}
