package usecase

import (
	"context"
	"log"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

const (
	defaultSettlementInterval   = 10 * time.Minute
	defaultAutoReceiveAfter     = 48 * time.Hour
	defaultSettlementHoldPeriod = 24 * time.Hour
	defaultSettlementBatchSize  = 100
)

// OrderSettlementScheduler runs periodic tasks to auto-receive and settle orders.
type OrderSettlementScheduler struct {
	orderRepo        domain.OrderRepository
	interval         time.Duration
	autoReceiveAfter time.Duration
	settlementAfter  time.Duration
	batchSize        int
	clock            func() time.Time
}

// NewOrderSettlementScheduler creates a scheduler with default settings.
func NewOrderSettlementScheduler(orderRepo domain.OrderRepository) *OrderSettlementScheduler {
	return &OrderSettlementScheduler{
		orderRepo:        orderRepo,
		interval:         defaultSettlementInterval,
		autoReceiveAfter: defaultAutoReceiveAfter,
		settlementAfter:  defaultSettlementHoldPeriod,
		batchSize:        defaultSettlementBatchSize,
		clock:            time.Now,
	}
}

// Start runs the scheduler loop until the context is canceled.
func (s *OrderSettlementScheduler) Start(ctx context.Context) {
	log.Printf("[scheduler] order settlement started (interval=%s)", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[scheduler] order settlement stopped")
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *OrderSettlementScheduler) runOnce(ctx context.Context) {
	if err := s.processAutoReceive(ctx); err != nil {
		log.Printf("[scheduler] auto-receive failed: %v", err)
	}
	if err := s.processSettlement(ctx); err != nil {
		log.Printf("[scheduler] settlement failed: %v", err)
	}
}

func (s *OrderSettlementScheduler) processAutoReceive(ctx context.Context) error {
	cutoff := s.clock().Add(-s.autoReceiveAfter)
	orders, err := s.orderRepo.ListAutoReceiveCandidates(ctx, cutoff, s.batchSize)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}

	appliedCount := 0
	for _, order := range orders {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		notes := "Order auto-received after delivery confirmation timeout"
		applied, err := s.orderRepo.MarkOrderReceived(ctx, order.ID, nil, &notes)
		if err != nil {
			log.Printf("[scheduler] auto-receive order %s failed: %v", order.ID, err)
			continue
		}
		if applied {
			appliedCount++
		}
	}
	if appliedCount > 0 {
		log.Printf("[scheduler] auto-received %d orders", appliedCount)
	}
	return nil
}

func (s *OrderSettlementScheduler) processSettlement(ctx context.Context) error {
	cutoff := s.clock().Add(-s.settlementAfter)
	orders, err := s.orderRepo.ListSettlementCandidates(ctx, cutoff, s.batchSize)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}

	appliedCount := 0
	for _, order := range orders {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		notes := "Order settled after holding period"
		applied, err := s.orderRepo.MarkOrderCompleted(ctx, order.ID, nil, &notes)
		if err != nil {
			log.Printf("[scheduler] settle order %s failed: %v", order.ID, err)
			continue
		}
		if applied {
			appliedCount++
		}
	}
	if appliedCount > 0 {
		log.Printf("[scheduler] settled %d orders", appliedCount)
	}
	return nil
}
