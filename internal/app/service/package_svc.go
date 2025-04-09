package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"math"
	"net/http"
	"time"

	"go.uber.org/dig"
)

type (
	PackageSvc interface {
		CreatePackage(ctx context.Context, req models.Package) (resp models.DefaultResponse, err error)
		GetPackages(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		GetPackageByID(ctx context.Context, packageId int64) (resp models.DefaultResponse, err error)
		UpdatePackage(ctx context.Context, req models.Package) (resp models.DefaultResponse, err error)
		DeletePackage(ctx context.Context, packageId int64, deletedBy string) (resp models.DefaultResponse, err error)
		CheckPackageAvailability(ctx context.Context) (bool, error)
	}

	PackageSvcImpl struct {
		dig.In

		PackageRepo postgres.PackageRepo
		ConfigRepo  postgres.ConfigRepo
	}
)

func NewPackageSvc(impl PackageSvcImpl) PackageSvc {
	return &impl
}

func (p *PackageSvcImpl) CreatePackage(ctx context.Context, req models.Package) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusCreated
		resp.Message = "Success"
	}

	discountExpired := getDiscountExpired(time.Now())
	if req.DiscountedPrice >= 0 {
		err := p.PackageRepo.UpdatePackageDiscountExpired(ctx, discountExpired)
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = "bad request"
			resp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][CreatePackage] while UpdatePackageDiscountExpired err : %v", err.Error()))

			return resp, err
		}
	}

	req.DiscountExpired = discountExpired.Format(time.RFC3339)
	packageId, err := p.PackageRepo.CreatePackage(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CreatePackage] while CreatePackage err : %v", err.Error()))

		return resp, err
	}

	type Package struct {
		ID int `json:"id"`
	}

	resp.Data = Package{ID: packageId}

	return
}

func (p *PackageSvcImpl) GetPackages(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error) {
	var dataResp models.DefaultResponse
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "Success"
	}

	packages, totalCount, err := p.PackageRepo.GetPackages(ctx, req)
	if err != nil {
		dataResp.Code = http.StatusBadRequest
		dataResp.Message = "bad request"
		dataResp.Error = err.Error()
		resp.Results = dataResp
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackages] while GetPackages err : %v", err.Error()))

		return resp, err
	}

	isAvailable, _ := p.CheckPackageAvailability(ctx)
	// if err != nil {
	// 	slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackages] while CheckPackageAvailability err : %v", err.Error()))
	// 	isAvailable = false
	// }

	// Set nilai IsAvailable ke setiap package
	for i := range packages {
		packages[i].IsAvailable = isAvailable
	}

	// mapping response
	{
		dataResp.Data = packages
		resp.Page = uint(req.Page)
		resp.Limit = uint(req.Limit)

		totalPage := math.Ceil(float64(totalCount) / float64(req.Limit))
		resp.TotalPages = uint(totalPage)

		resp.TotalItems = uint(totalCount)
		resp.HasNext = req.Page < int64(resp.TotalPages)
		resp.HasPrevious = req.Page > 1
		resp.Results = dataResp
	}

	return
}

func (p *PackageSvcImpl) GetPackageByID(ctx context.Context, packageId int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	packageData, err := p.PackageRepo.GetPackageByID(ctx, packageId)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackageByID] while GetPackageByID err : %v", err))

		return resp, err
	}

	isAvailable, _ := p.CheckPackageAvailability(ctx)
	// if err != nil {
	// 	slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackageByID] while CheckPackageAvailability err : %v", err.Error()))
	// 	isAvailable = false
	// }

	packageData.IsAvailable = isAvailable

	resp.Data = packageData

	return
}

func (p *PackageSvcImpl) UpdatePackage(ctx context.Context, req models.Package) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	packageData, err := p.PackageRepo.GetPackageByID(ctx, req.ID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] while GetPackageByID err : %v", err))

		return resp, err
	}

	if packageData.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = fmt.Sprintf("package with id %d not found", req.ID)
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] package with id %d not found", req.ID))

		return resp, fmt.Errorf("package with id %d not found", req.ID)
	}

	if req.DiscountedPrice != 0 && req.DiscountedPrice != packageData.DiscountedPrice {
		newDiscountExpired := getDiscountExpired(time.Now())
		err = p.PackageRepo.UpdatePackageDiscountExpired(ctx, newDiscountExpired)
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = "bad request"
			resp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] while UpdatePackageDiscountExpired err : %v", err))

			return resp, err
		}
	}

	err = p.PackageRepo.UpdatePackage(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] while UpdatePackage err : %v", err))

		return resp, err
	}

	return
}

func (p *PackageSvcImpl) DeletePackage(ctx context.Context, packageId int64, deletedBy string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	packageData, err := p.PackageRepo.GetPackageByID(ctx, packageId)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeletePackage] while GetPackageByID err : %v", err))

		return resp, err
	}

	if packageData.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = fmt.Sprintf("package with id %d not found", packageId)
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeletePackage] package with id %d not found", packageId))

		return resp, fmt.Errorf("package with id %d not found", packageId)
	}

	err = p.PackageRepo.SoftDeletePackage(ctx, packageId, deletedBy)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeletePackage] while DeletePackage err : %v", err))

		return resp, err
	}

	return
}

func getDiscountExpired(t time.Time) time.Time {
	firstDayNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstDayNextMonth.Add(-time.Second)
}

func (p *PackageSvcImpl) CheckPackageAvailability(ctx context.Context) (bool, error) {
	startDate, err := p.ConfigRepo.GetConfigByKey(utils.MembershipProgramStartDateConfig)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CheckPackageAvailability] while GetConfigByKey %s err : %v", utils.MembershipProgramStartDateConfig, err))
		return false, fmt.Errorf("failed to get %s: %w", utils.MembershipProgramStartDateConfig, err)
	}

	endDate, err := p.ConfigRepo.GetConfigByKey(utils.MembershipProgramEndDateConfig)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CheckPackageAvailability] while GetConfigByKey %s err : %v", utils.MembershipProgramEndDateConfig, err))
		return false, fmt.Errorf("failed to get %s: %w", utils.MembershipProgramEndDateConfig, err)
	}

	layout := "2006-01-02"

	startTime, err := time.Parse(layout, startDate)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CheckPackageAvailability] while Parse %s err : %v", utils.MembershipProgramStartDateConfig, err))
		return false, fmt.Errorf("failed to parse %s: %w", utils.MembershipProgramStartDateConfig, err)
	}

	endTime, err := time.Parse(layout, endDate)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CheckPackageAvailability] while Parse %s err : %v", utils.MembershipProgramEndDateConfig, err))
		return false, fmt.Errorf("failed to parse %s: %w", utils.MembershipProgramEndDateConfig, err)
	}

	currentTime := time.Now()

	// Compare berdasarkan tanggal saja
	currentDate := currentTime.Truncate(24 * time.Hour)
	if currentDate.Before(startTime) || currentDate.After(endTime) {
		return true, nil
	}

	return false, nil
}
