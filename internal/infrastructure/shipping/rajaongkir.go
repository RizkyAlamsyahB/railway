package shipping

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// cacheTTL for static location data (provinces, cities, districts, subdistricts).
// Only successful responses are cached; errors are never cached.
const cacheTTL = 365 * 24 * time.Hour // 12 months

type cacheEntry struct {
	data      json.RawMessage
	expiresAt time.Time
}

// rajaOngkirClient implements domain.RajaOngkirProvider.
type rajaOngkirClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	mu         sync.RWMutex
	cache      map[string]cacheEntry
}

// NewRajaOngkirClient creates a new RajaOngkir API client.
func NewRajaOngkirClient(cfg config.RajaOngkirConfig) (domain.RajaOngkirProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("RAJAONGKIR_API_KEY is required")
	}
	return &rajaOngkirClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		cache: make(map[string]cacheEntry),
	}, nil
}

// komerceEnvelope is the generic envelope from Komerce RajaOngkir API.
// Response shape: {"meta": {"message":"...", "code":200, "status":"success"}, "data": [...]}
type komerceEnvelope struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data json.RawMessage `json:"data"`
}

func (c *rajaOngkirClient) doGet(ctx context.Context, path string) (json.RawMessage, error) {
	// Check cache first
	c.mu.RLock()
	if entry, ok := c.cache[path]; ok && time.Now().Before(entry.expiresAt) {
		c.mu.RUnlock()
		return entry.data, nil
	}
	c.mu.RUnlock()

	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to create request: %w", err)
	}
	req.Header.Set("key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to read response body: %w", err)
	}

	// On non-200 status, return error WITHOUT caching
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rajaongkir: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Komerce envelope: {"meta":{...}, "data":[...]}
	var envelope komerceEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Data != nil {
		// Cache successful response only
		c.mu.Lock()
		c.cache[path] = cacheEntry{data: envelope.Data, expiresAt: time.Now().Add(cacheTTL)}
		c.mu.Unlock()
		return envelope.Data, nil
	}

	// Fallback: return raw body as JSON (also cache)
	c.mu.Lock()
	c.cache[path] = cacheEntry{data: body, expiresAt: time.Now().Add(cacheTTL)}
	c.mu.Unlock()
	return body, nil
}

// locationItem is the common shape returned by Komerce location endpoints.
// Province/City/District: {"id": <int>, "name": "<string>"}
// SubDistrict also includes: {"zip_code": "<string>"}
type locationItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	ZipCode string `json:"zip_code,omitempty"`
}

// GetProvinces returns all provinces from RajaOngkir.
// Endpoint: GET /destination/province
func (c *rajaOngkirClient) GetProvinces(ctx context.Context) ([]domain.ROProvince, error) {
	data, err := c.doGet(ctx, "/destination/province")
	if err != nil {
		return nil, err
	}

	var items []locationItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to parse provinces: %w", err)
	}

	provinces := make([]domain.ROProvince, len(items))
	for i, item := range items {
		provinces[i] = domain.ROProvince{
			ID:   strconv.Itoa(item.ID),
			Name: item.Name,
		}
	}
	return provinces, nil
}

// GetCitiesByProvince returns cities in a province.
// Endpoint: GET /destination/city/{province_id}
func (c *rajaOngkirClient) GetCitiesByProvince(ctx context.Context, provinceID string) ([]domain.ROCity, error) {
	data, err := c.doGet(ctx, "/destination/city/"+provinceID)
	if err != nil {
		return nil, err
	}

	var items []locationItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to parse cities: %w", err)
	}

	cities := make([]domain.ROCity, len(items))
	for i, item := range items {
		cities[i] = domain.ROCity{
			ID:         strconv.Itoa(item.ID),
			ProvinceID: provinceID,
			Name:       item.Name,
		}
	}
	return cities, nil
}

// GetDistrictsByCity returns districts in a city.
// Endpoint: GET /destination/district/{city_id}
func (c *rajaOngkirClient) GetDistrictsByCity(ctx context.Context, cityID string) ([]domain.RODistrict, error) {
	data, err := c.doGet(ctx, "/destination/district/"+cityID)
	if err != nil {
		return nil, err
	}

	var items []locationItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to parse districts: %w", err)
	}

	districts := make([]domain.RODistrict, len(items))
	for i, item := range items {
		districts[i] = domain.RODistrict{
			ID:     strconv.Itoa(item.ID),
			CityID: cityID,
			Name:   item.Name,
		}
	}
	return districts, nil
}

// GetSubdistrictsByDistrict returns subdistricts in a district.
// Endpoint: GET /destination/sub-district/{district_id}
func (c *rajaOngkirClient) GetSubdistrictsByDistrict(ctx context.Context, districtID string) ([]domain.ROSubdistrict, error) {
	data, err := c.doGet(ctx, "/destination/sub-district/"+districtID)
	if err != nil {
		return nil, err
	}

	var items []locationItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("rajaongkir: failed to parse subdistricts: %w", err)
	}

	subdistricts := make([]domain.ROSubdistrict, len(items))
	for i, item := range items {
		subdistricts[i] = domain.ROSubdistrict{
			ID:         strconv.Itoa(item.ID),
			DistrictID: districtID,
			Name:       item.Name,
			ZipCode:    item.ZipCode,
		}
	}
	return subdistricts, nil
}

// costItem is the shape of each item in the RajaOngkir domestic cost response.
type costItem struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	ETD         string `json:"etd"`
}

// CalculateDomesticCost calls RajaOngkir Komerce domestic cost API.
// POST /calculate/district/domestic-cost
// Form data: origin, destination (district IDs), weight (grams), courier (colon-separated), price=lowest
func (c *rajaOngkirClient) CalculateDomesticCost(ctx context.Context, originDistrictID, destDistrictID string, weightGram int, courierCodes string) ([]domain.ShippingCostOption, error) {
	apiURL := c.baseURL + "/calculate/district/domestic-cost"

	form := url.Values{}
	form.Set("origin", originDistrictID)
	form.Set("destination", destDistrictID)
	form.Set("weight", strconv.Itoa(weightGram))
	form.Set("courier", courierCodes)
	form.Set("price", "lowest")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("rajaongkir cost: failed to create request: %w", err)
	}
	req.Header.Set("key", c.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir cost: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir cost: failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rajaongkir cost: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var envelope komerceEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("rajaongkir cost: failed to parse envelope: %w", err)
	}

	if envelope.Data == nil {
		return nil, fmt.Errorf("rajaongkir cost: empty data in response")
	}

	var items []costItem
	if err := json.Unmarshal(envelope.Data, &items); err != nil {
		return nil, fmt.Errorf("rajaongkir cost: failed to parse cost items: %w", err)
	}

	options := make([]domain.ShippingCostOption, len(items))
	for i, item := range items {
		options[i] = domain.ShippingCostOption{
			Name:        item.Name,
			Code:        item.Code,
			Service:     item.Service,
			Description: item.Description,
			Cost:        item.Cost,
			ETD:         item.ETD,
		}
	}
	return options, nil
}

// TrackWaybill tracks a shipment via RajaOngkir Komerce waybill API.
// POST /track/waybill?awb={awb}&courier={courier}
func (c *rajaOngkirClient) TrackWaybill(ctx context.Context, awbNumber, courierCode string) (*domain.TrackWaybillResponse, error) {
	apiURL := c.baseURL + "/track/waybill?awb=" + url.QueryEscape(awbNumber) + "&courier=" + url.QueryEscape(courierCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir track: failed to create request: %w", err)
	}
	req.Header.Set("key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir track: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rajaongkir track: failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rajaongkir track: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the Komerce envelope.
	var envelope komerceEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("rajaongkir track: failed to parse envelope: %w", err)
	}
	if envelope.Data == nil {
		return nil, fmt.Errorf("rajaongkir track: empty data in response")
	}

	// Parse the waybill tracking data.
	var waybillResp waybillAPIResponse
	if err := json.Unmarshal(envelope.Data, &waybillResp); err != nil {
		return nil, fmt.Errorf("rajaongkir track: failed to parse waybill data: %w", err)
	}

	return mapWaybillResponse(&waybillResp), nil
}

// --- waybill API response structs ---

type waybillAPIResponse struct {
	Delivered      bool                  `json:"delivered"`
	Summary        waybillSummary        `json:"summary"`
	Details        waybillDetails        `json:"details"`
	DeliveryStatus waybillDeliveryStatus `json:"delivery_status"`
	Manifest       []waybillManifest     `json:"manifest"`
}

type waybillSummary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

type waybillDetails struct {
	WaybillNumber    string `json:"waybill_number"`
	WaybillDate      string `json:"waybill_date"`
	WaybillTime      string `json:"waybill_time"`
	Weight           string `json:"weight"`
	Origin           string `json:"origin"`
	Destination      string `json:"destination"`
	ShipperName      string `json:"shipper_name"`
	ShipperAddress1  string `json:"shipper_address1"`
	ShipperAddress2  string `json:"shipper_address2"`
	ShipperAddress3  string `json:"shipper_address3"`
	ShipperCity      string `json:"shipper_city"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverAddress1 string `json:"receiver_address1"`
	ReceiverAddress2 string `json:"receiver_address2"`
	ReceiverAddress3 string `json:"receiver_address3"`
	ReceiverCity     string `json:"receiver_city"`
}

type waybillDeliveryStatus struct {
	Status      string `json:"status"`
	PodReceiver string `json:"pod_receiver"`
	PodDate     string `json:"pod_date"`
	PodTime     string `json:"pod_time"`
}

type waybillManifest struct {
	ManifestCode        string `json:"manifest_code"`
	ManifestDescription string `json:"manifest_description"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name"`
}

func mapWaybillResponse(r *waybillAPIResponse) *domain.TrackWaybillResponse {
	manifest := make([]domain.TrackManifestItem, len(r.Manifest))
	for i, m := range r.Manifest {
		manifest[i] = domain.TrackManifestItem{
			ManifestCode:        m.ManifestCode,
			ManifestDescription: m.ManifestDescription,
			ManifestDate:        m.ManifestDate,
			ManifestTime:        m.ManifestTime,
			CityName:            m.CityName,
		}
	}
	return &domain.TrackWaybillResponse{
		Delivered: r.Delivered,
		Summary: domain.TrackWaybillSummary{
			CourierCode:   r.Summary.CourierCode,
			CourierName:   r.Summary.CourierName,
			WaybillNumber: r.Summary.WaybillNumber,
			ServiceCode:   r.Summary.ServiceCode,
			WaybillDate:   r.Summary.WaybillDate,
			ShipperName:   r.Summary.ShipperName,
			ReceiverName:  r.Summary.ReceiverName,
			Origin:        r.Summary.Origin,
			Destination:   r.Summary.Destination,
			Status:        r.Summary.Status,
		},
		Details: domain.TrackWaybillDetails{
			WaybillNumber:    r.Details.WaybillNumber,
			WaybillDate:      r.Details.WaybillDate,
			WaybillTime:      r.Details.WaybillTime,
			Weight:           r.Details.Weight,
			Origin:           r.Details.Origin,
			Destination:      r.Details.Destination,
			ShipperName:      r.Details.ShipperName,
			ShipperAddress1:  r.Details.ShipperAddress1,
			ShipperAddress2:  r.Details.ShipperAddress2,
			ShipperAddress3:  r.Details.ShipperAddress3,
			ShipperCity:      r.Details.ShipperCity,
			ReceiverName:     r.Details.ReceiverName,
			ReceiverAddress1: r.Details.ReceiverAddress1,
			ReceiverAddress2: r.Details.ReceiverAddress2,
			ReceiverAddress3: r.Details.ReceiverAddress3,
			ReceiverCity:     r.Details.ReceiverCity,
		},
		DeliveryStatus: domain.TrackDeliveryStatus{
			Status:      r.DeliveryStatus.Status,
			PodReceiver: r.DeliveryStatus.PodReceiver,
			PodDate:     r.DeliveryStatus.PodDate,
			PodTime:     r.DeliveryStatus.PodTime,
		},
		Manifest: manifest,
	}
}
