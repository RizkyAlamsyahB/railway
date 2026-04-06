package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
)

func TestToVendorModel_AllowsNilVendorType(t *testing.T) {
	now := time.Now()
	vendor := &domain.Vendor{
		ID:          uuid.New(),
		OwnerUserID: uuid.New(),
		VendorType:  nil,
		DisplayName: nil,
		Status:      domain.VendorStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	model := toVendorModel(vendor)
	if model.VendorType != nil {
		t.Fatalf("expected nil vendor type in model, got %v", model.VendorType)
	}
}

func TestToDomainVendor_AllowsNilVendorType(t *testing.T) {
	now := time.Now()
	model := &vendorModel{
		ID:          uuid.NewString(),
		OwnerUserID: uuid.NewString(),
		VendorType:  nil,
		DisplayName: nil,
		Status:      domain.VendorStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	vendor := toDomainVendor(model)
	if vendor.VendorType != nil {
		t.Fatalf("expected nil vendor type in domain, got %v", vendor.VendorType)
	}
}

func TestToDomainVendor_PreservesVendorTypeValue(t *testing.T) {
	now := time.Now()
	vendorType := domain.VendorTypePPIU
	displayName := "Vendor PPIU"
	model := &vendorModel{
		ID:          uuid.NewString(),
		OwnerUserID: uuid.NewString(),
		VendorType:  &vendorType,
		DisplayName: &displayName,
		Status:      domain.VendorStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	vendor := toDomainVendor(model)
	if vendor.VendorType == nil || *vendor.VendorType != domain.VendorTypePPIU {
		t.Fatalf("expected vendor type %q, got %v", domain.VendorTypePPIU, vendor.VendorType)
	}
}

func TestVendorModelMapping_PreservesStructuredAddress(t *testing.T) {
	now := time.Now()
	provinceID := "31"
	provinceName := "DKI Jakarta"
	cityID := "3171"
	cityName := "Jakarta Pusat"
	districtID := "317101"
	districtName := "Menteng"
	subdistrictID := "3171011001"
	subdistrictName := "Pegangsaan"
	postalCode := "10320"
	addressLine := "Jl. Pegangsaan Barat No. 12"

	vendor := &domain.Vendor{
		ID:              uuid.New(),
		OwnerUserID:     uuid.New(),
		ProvinceID:      &provinceID,
		ProvinceName:    &provinceName,
		CityID:          &cityID,
		CityName:        &cityName,
		DistrictID:      &districtID,
		DistrictName:    &districtName,
		SubdistrictID:   &subdistrictID,
		SubdistrictName: &subdistrictName,
		PostalCode:      &postalCode,
		AddressLine:     &addressLine,
		Status:          domain.VendorStatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	model := toVendorModel(vendor)
	if model.ProvinceID == nil || *model.ProvinceID != provinceID {
		t.Fatalf("expected province_id %q, got %v", provinceID, model.ProvinceID)
	}
	if model.AddressLine == nil || *model.AddressLine != addressLine {
		t.Fatalf("expected address_line %q, got %v", addressLine, model.AddressLine)
	}

	roundTrip := toDomainVendor(&model)
	if roundTrip.SubdistrictName == nil || *roundTrip.SubdistrictName != subdistrictName {
		t.Fatalf("expected subdistrict_name %q, got %v", subdistrictName, roundTrip.SubdistrictName)
	}
	if roundTrip.PostalCode == nil || *roundTrip.PostalCode != postalCode {
		t.Fatalf("expected postal_code %q, got %v", postalCode, roundTrip.PostalCode)
	}
}

func TestBuildVendorProposalUpdates_IncludesStructuredAddressNames(t *testing.T) {
	now := time.Now()
	provinceID := "31"
	provinceName := "DKI Jakarta"
	cityID := "3171"
	cityName := "Jakarta Pusat"
	districtID := "317101"
	districtName := "Menteng"
	subdistrictID := "3171011001"
	subdistrictName := "Pegangsaan"
	postalCode := "10320"
	addressLine := "Jl. Pegangsaan Barat No. 12"

	updates := buildVendorProposalUpdates(&domain.Vendor{
		VendorType:      ptrString(domain.VendorTypeSouvenirStore),
		DisplayName:     ptrString("Toko Haji"),
		Description:     ptrString("Pusat oleh-oleh"),
		ProvinceID:      &provinceID,
		ProvinceName:    &provinceName,
		CityID:          &cityID,
		CityName:        &cityName,
		DistrictID:      &districtID,
		DistrictName:    &districtName,
		SubdistrictID:   &subdistrictID,
		SubdistrictName: &subdistrictName,
		PostalCode:      &postalCode,
		AddressLine:     &addressLine,
		Status:          domain.VendorStatusSubmitted,
		UpdatedAt:       now,
	})

	if got, ok := updates["province_name"].(*string); !ok || got == nil || *got != provinceName {
		t.Fatalf("expected province_name %q, got %#v", provinceName, updates["province_name"])
	}
	if got, ok := updates["city_name"].(*string); !ok || got == nil || *got != cityName {
		t.Fatalf("expected city_name %q, got %#v", cityName, updates["city_name"])
	}
	if got, ok := updates["district_name"].(*string); !ok || got == nil || *got != districtName {
		t.Fatalf("expected district_name %q, got %#v", districtName, updates["district_name"])
	}
	if got, ok := updates["subdistrict_name"].(*string); !ok || got == nil || *got != subdistrictName {
		t.Fatalf("expected subdistrict_name %q, got %#v", subdistrictName, updates["subdistrict_name"])
	}
}

func ptrString(v string) *string {
	return &v
}

func TestVendorResponsiblePersonModel_EncryptsAndDecryptsNIK(t *testing.T) {
	cipher, err := sensitivedata.NewFieldCipher([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("failed to create field cipher: %v", err)
	}

	repo := &vendorRepository{fieldCipher: cipher}
	now := time.Now()
	rp := &domain.VendorResponsiblePerson{
		ID:        uuid.New(),
		VendorID:  uuid.New(),
		UserID:    uuid.New(),
		NIK:       "3173010101010001",
		CreatedAt: now,
		UpdatedAt: now,
	}

	model, err := repo.toVendorResponsiblePersonModel(rp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if model.NIK == rp.NIK {
		t.Fatalf("expected encrypted NIK, got plaintext %q", model.NIK)
	}

	roundTrip, err := repo.toDomainVendorResponsiblePerson(&model)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if roundTrip.NIK != rp.NIK {
		t.Fatalf("expected decrypted NIK %q, got %q", rp.NIK, roundTrip.NIK)
	}
}
