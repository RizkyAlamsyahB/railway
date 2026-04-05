package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func setupAddressUseCase(t *testing.T) (*mocks.MockAddressRepository, domain.AddressUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockAddressRepository(ctrl)
	uc := NewAddressUseCase(repo)
	return repo, uc
}

// ============================================================
// Delete
// ============================================================

func TestDelete_Success(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()
	defaultAddrID := uuid.New()

	// Non-default address, user has 2 addresses
	addr := &domain.Address{ID: addrID, UserID: userID, IsDefault: false}
	repo.EXPECT().FindByID(ctx, addrID).Return(addr, nil)
	repo.EXPECT().FindByUserID(ctx, userID).Return([]domain.Address{
		{ID: defaultAddrID, UserID: userID, IsDefault: true},
		{ID: addrID, UserID: userID, IsDefault: false},
	}, nil)
	repo.EXPECT().Delete(ctx, addrID).Return(nil)

	err := uc.Delete(ctx, userID, addrID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_AddressNotFound(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()

	repo.EXPECT().FindByID(ctx, addrID).Return(nil, nil)

	err := uc.Delete(ctx, userID, addrID)
	if !errors.Is(err, ErrAddressNotFound) {
		t.Errorf("expected ErrAddressNotFound, got %v", err)
	}
}

func TestDelete_CannotDeleteDefaultAddress(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()

	// Default address
	addr := &domain.Address{ID: addrID, UserID: userID, IsDefault: true}
	repo.EXPECT().FindByID(ctx, addrID).Return(addr, nil)

	err := uc.Delete(ctx, userID, addrID)
	if !errors.Is(err, ErrCannotDeleteDefaultAddr) {
		t.Errorf("expected ErrCannotDeleteDefaultAddr, got %v", err)
	}
}

func TestDelete_CannotDeleteLastAddress(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()

	// Non-default, but only 1 address
	addr := &domain.Address{ID: addrID, UserID: userID, IsDefault: false}
	repo.EXPECT().FindByID(ctx, addrID).Return(addr, nil)
	repo.EXPECT().FindByUserID(ctx, userID).Return([]domain.Address{
		{ID: addrID, UserID: userID, IsDefault: false},
	}, nil)

	err := uc.Delete(ctx, userID, addrID)
	if !errors.Is(err, ErrCannotDeleteLastAddr) {
		t.Errorf("expected ErrCannotDeleteLastAddr, got %v", err)
	}
}

func TestDelete_FindByIDError(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()

	repo.EXPECT().FindByID(ctx, addrID).Return(nil, errors.New("db error"))

	err := uc.Delete(ctx, userID, addrID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_FindByUserIDError(t *testing.T) {
	repo, uc := setupAddressUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	addrID := uuid.New()

	addr := &domain.Address{ID: addrID, UserID: userID, IsDefault: false}
	repo.EXPECT().FindByID(ctx, addrID).Return(addr, nil)
	repo.EXPECT().FindByUserID(ctx, userID).Return(nil, errors.New("db error"))

	err := uc.Delete(ctx, userID, addrID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
