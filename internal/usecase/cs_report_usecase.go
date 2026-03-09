package usecase

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// ============================================================
// CS Report UseCase
// ============================================================

type csReportUseCase struct {
	ticketRepo domain.TicketRepository
}

func NewCSReportUseCase(ticketRepo domain.TicketRepository) domain.CSReportUseCase {
	return &csReportUseCase{ticketRepo: ticketRepo}
}

func parseMonth(month string) (int, int, error) {
	parts := strings.SplitN(month, "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid month format")
	}
	y, err := strconv.Atoi(parts[0])
	if err != nil || y < 2000 || y > 2100 {
		return 0, 0, fmt.Errorf("invalid year")
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 1 || m > 12 {
		return 0, 0, fmt.Errorf("invalid month")
	}
	return y, m, nil
}

func (uc *csReportUseCase) GetReport(ctx context.Context, params domain.CSReportParams) (*domain.CSReportResponse, error) {
	y, m, err := parseMonth(params.Month)
	if err != nil {
		return nil, ErrInvalidMonth
	}

	total, closed, err := uc.ticketRepo.CountByMonth(ctx, y, m)
	if err != nil {
		return nil, fmt.Errorf("failed to count tickets: %w", err)
	}

	avgMin, err := uc.ticketRepo.AvgFirstResponseMinute(ctx, y, m)
	if err != nil {
		return nil, fmt.Errorf("failed to get avg response time: %w", err)
	}
	// Round to 1 decimal
	avgMin = math.Round(avgMin*10) / 10

	topSubjects, err := uc.ticketRepo.TopSubjectsByMonth(ctx, y, m, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top subjects: %w", err)
	}

	perDay, err := uc.ticketRepo.TicketsPerDayByMonth(ctx, y, m)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets per day: %w", err)
	}

	return &domain.CSReportResponse{
		Month:             params.Month,
		TotalTickets:      total,
		ClosedTickets:     closed,
		AvgResponseMinute: avgMin,
		TopSubjects:       topSubjects,
		TicketsPerDay:     perDay,
	}, nil
}

func (uc *csReportUseCase) ExportReport(ctx context.Context, params domain.CSReportParams) ([]domain.CSReportExportRow, error) {
	y, m, err := parseMonth(params.Month)
	if err != nil {
		return nil, ErrInvalidMonth
	}
	return uc.ticketRepo.ExportByMonth(ctx, y, m)
}

// ============================================================
// Ticket Subject UseCase
// ============================================================

type ticketSubjectUseCase struct {
	subjectRepo domain.TicketSubjectRepository
}

func NewTicketSubjectUseCase(subjectRepo domain.TicketSubjectRepository) domain.TicketSubjectUseCase {
	return &ticketSubjectUseCase{subjectRepo: subjectRepo}
}

func (uc *ticketSubjectUseCase) ListSubjects(ctx context.Context) ([]domain.TicketSubject, error) {
	return uc.subjectRepo.ListActive(ctx)
}
