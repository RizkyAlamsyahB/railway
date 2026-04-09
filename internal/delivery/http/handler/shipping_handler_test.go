package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	authutil "github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

const shippingHandlerTestJWTSecret = "shipping-handler-test-secret"

type stubShippingUseCase struct {
	getProvinces    func(context.Context) ([]domain.ROProvince, error)
	getCities       func(context.Context, string) ([]domain.ROCity, error)
	getDistricts    func(context.Context, string) ([]domain.RODistrict, error)
	getSubdistricts func(context.Context, string) ([]domain.ROSubdistrict, error)
}

func (s stubShippingUseCase) GetProvinces(ctx context.Context) ([]domain.ROProvince, error) {
	return s.getProvinces(ctx)
}

func (s stubShippingUseCase) GetCities(ctx context.Context, provinceID string) ([]domain.ROCity, error) {
	return s.getCities(ctx, provinceID)
}

func (s stubShippingUseCase) GetDistricts(ctx context.Context, cityID string) ([]domain.RODistrict, error) {
	return s.getDistricts(ctx, cityID)
}

func (s stubShippingUseCase) GetSubdistricts(ctx context.Context, districtID string) ([]domain.ROSubdistrict, error) {
	return s.getSubdistricts(ctx, districtID)
}

func TestShippingHandlerLocationAliasParity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewShippingHandler(stubShippingUseCase{
		getProvinces: func(context.Context) ([]domain.ROProvince, error) {
			return []domain.ROProvince{
				{ID: "1", Name: "DKI Jakarta"},
				{ID: "2", Name: "Jawa Barat"},
			}, nil
		},
		getCities: func(_ context.Context, provinceID string) ([]domain.ROCity, error) {
			return []domain.ROCity{
				{ID: "10", ProvinceID: provinceID, Name: "Jakarta Selatan", Type: "Kota", PostalCode: "12110"},
			}, nil
		},
		getDistricts: func(_ context.Context, cityID string) ([]domain.RODistrict, error) {
			return []domain.RODistrict{
				{ID: "100", CityID: cityID, Name: "Kebayoran Baru"},
			}, nil
		},
		getSubdistricts: func(_ context.Context, districtID string) ([]domain.ROSubdistrict, error) {
			return []domain.ROSubdistrict{
				{ID: "1000", DistrictID: districtID, Name: "Senayan", ZipCode: "12190"},
			}, nil
		},
	})

	r := newShippingAliasTestRouter(h)
	token := mustGenerateTestToken(t, "customer")

	testCases := []struct {
		name string
		path string
	}{
		{name: "provinces", path: "/provinces"},
		{name: "cities", path: "/cities?province_id=6"},
		{name: "districts", path: "/districts?city_id=39"},
		{name: "subdistricts", path: "/subdistricts?district_id=512"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shippingResp := performRequest(r, http.MethodGet, "/api/v1/shipping"+tc.path, token)
			locationResp := performRequest(r, http.MethodGet, "/api/v1/locations"+tc.path, token)

			if shippingResp.Code != http.StatusOK {
				t.Fatalf("shipping status = %d, want 200", shippingResp.Code)
			}
			if locationResp.Code != http.StatusOK {
				t.Fatalf("locations status = %d, want 200", locationResp.Code)
			}
			if shippingResp.Body.String() != locationResp.Body.String() {
				t.Fatalf("expected identical responses, shipping=%s locations=%s", shippingResp.Body.String(), locationResp.Body.String())
			}
		})
	}
}

func TestShippingHandlerLocationAliasValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewShippingHandler(stubShippingUseCase{
		getProvinces: func(context.Context) ([]domain.ROProvince, error) { return nil, nil },
		getCities: func(context.Context, string) ([]domain.ROCity, error) {
			t.Fatal("GetCities should not be called when province_id is missing")
			return nil, nil
		},
		getDistricts: func(context.Context, string) ([]domain.RODistrict, error) {
			t.Fatal("GetDistricts should not be called when city_id is missing")
			return nil, nil
		},
		getSubdistricts: func(context.Context, string) ([]domain.ROSubdistrict, error) {
			t.Fatal("GetSubdistricts should not be called when district_id is missing")
			return nil, nil
		},
	})

	r := newShippingAliasTestRouter(h)
	token := mustGenerateTestToken(t, "customer")

	testCases := []struct {
		name    string
		path    string
		message string
	}{
		{name: "cities requires province", path: "/api/v1/locations/cities", message: "province_id is required"},
		{name: "districts requires city", path: "/api/v1/locations/districts", message: "city_id is required"},
		{name: "subdistricts requires district", path: "/api/v1/locations/subdistricts", message: "district_id is required"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := performRequest(r, http.MethodGet, tc.path, token)
			if resp.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", resp.Code)
			}

			var env responseEnvelope
			if err := json.Unmarshal(resp.Body.Bytes(), &env); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if env.Message != tc.message {
				t.Fatalf("message = %q, want %q", env.Message, tc.message)
			}
		})
	}
}

func TestShippingHandlerLocationAliasAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewShippingHandler(stubShippingUseCase{
		getProvinces: func(context.Context) ([]domain.ROProvince, error) {
			return []domain.ROProvince{{ID: "1", Name: "DKI Jakarta"}}, nil
		},
		getCities:       func(context.Context, string) ([]domain.ROCity, error) { return nil, nil },
		getDistricts:    func(context.Context, string) ([]domain.RODistrict, error) { return nil, nil },
		getSubdistricts: func(context.Context, string) ([]domain.ROSubdistrict, error) { return nil, nil },
	})

	r := newShippingAliasTestRouter(h)

	t.Run("customer allowed", func(t *testing.T) {
		resp := performRequest(r, http.MethodGet, "/api/v1/locations/provinces", mustGenerateTestToken(t, "customer"))
		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.Code)
		}
	})

	t.Run("vendor allowed", func(t *testing.T) {
		resp := performRequest(r, http.MethodGet, "/api/v1/locations/provinces", mustGenerateTestToken(t, "umkm"))
		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.Code)
		}
	})

	t.Run("other role forbidden", func(t *testing.T) {
		resp := performRequest(r, http.MethodGet, "/api/v1/locations/provinces", mustGenerateTestToken(t, "admin"))
		if resp.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", resp.Code)
		}

		var env responseEnvelope
		if err := json.Unmarshal(resp.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "insufficient permissions" {
			t.Fatalf("message = %q, want %q", env.Message, "insufficient permissions")
		}
	})

	t.Run("missing token unauthorized", func(t *testing.T) {
		resp := performRequest(r, http.MethodGet, "/api/v1/locations/provinces", "")
		if resp.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", resp.Code)
		}

		var env responseEnvelope
		if err := json.Unmarshal(resp.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "missing authorization header" {
			t.Fatalf("message = %q, want %q", env.Message, "missing authorization header")
		}
	})
}

func newShippingAliasTestRouter(h *ShippingHandler) *gin.Engine {
	r := gin.New()
	shippingGroup := r.Group("/api/v1/shipping")
	shippingGroup.Use(middleware.Auth(shippingHandlerTestJWTSecret))
	shippingGroup.Use(middleware.RequireRoles("customer", "umkm"))
	{
		shippingGroup.GET("/provinces", h.GetProvinces)
		shippingGroup.GET("/cities", h.GetCities)
		shippingGroup.GET("/districts", h.GetDistricts)
		shippingGroup.GET("/subdistricts", h.GetSubdistricts)
	}

	locationsGroup := r.Group("/api/v1/locations")
	locationsGroup.Use(middleware.Auth(shippingHandlerTestJWTSecret))
	locationsGroup.Use(middleware.RequireRoles("customer", "umkm"))
	{
		locationsGroup.GET("/provinces", h.GetProvinces)
		locationsGroup.GET("/cities", h.GetCities)
		locationsGroup.GET("/districts", h.GetDistricts)
		locationsGroup.GET("/subdistricts", h.GetSubdistricts)
	}

	return r
}

func performRequest(r *gin.Engine, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func mustGenerateTestToken(t *testing.T, role string) string {
	t.Helper()

	userID := uuid.New()
	var vendorID *uuid.UUID
	if role == "umkm" {
		id := uuid.New()
		vendorID = &id
	}

	token, err := authutil.GenerateToken(
		userID,
		"test@example.com",
		role,
		vendorID,
		nil,
		shippingHandlerTestJWTSecret,
		1,
		"shipping-handler-test",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}
