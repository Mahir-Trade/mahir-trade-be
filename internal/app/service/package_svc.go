package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
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
		CheckPackageAvailability(ctx context.Context, membershipPackage models.Package, userMembership models.UserMembership) (bool, error)
	}

	PackageSvcImpl struct {
		dig.In

		PackageRepo        postgres.PackageRepo
		ConfigRepo         postgres.ConfigRepo
		UserMembershipRepo postgres.UserMembershipRepo
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
			resp.Message = http.StatusText(http.StatusBadRequest)
			resp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][CreatePackage] while UpdatePackageDiscountExpired err : %v", err.Error()))

			return resp, err
		}
	}

	req.DiscountExpired = discountExpired.Format(time.RFC3339)
	packageId, err := p.PackageRepo.CreatePackage(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
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
		dataResp.Message = http.StatusText(http.StatusBadRequest)
		dataResp.Error = err.Error()
		resp.Results = dataResp
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackages] while GetPackages err : %v", err.Error()))

		return resp, err
	}

	userMembership := models.UserMembership{}
	if ctx.Value(middleware.UserData) != nil {
		userMembership, err = p.UserMembershipRepo.GetUserMembershipByUserID(ctx, ctx.Value(middleware.UserData).(middleware.UserCtxReq).UserID)
		if err != nil {
			dataResp.Code = http.StatusBadRequest
			dataResp.Message = http.StatusText(http.StatusBadRequest)
			dataResp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackageByID] while GetUserMembershipByUserID err : %v", err))

			return resp, err
		}
	}

	for i := range packages {
		packages[i].IsAvailable, _ = p.CheckPackageAvailability(ctx, packages[i], userMembership)
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
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackageByID] while GetPackageByID err : %v", err))

		return resp, err
	}

	userMembership, err := p.UserMembershipRepo.GetUserMembershipByUserID(ctx, ctx.Value(middleware.UserData).(middleware.UserCtxReq).UserID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetPackageByID] while GetUserMembershipByUserID err : %v", err))

		return resp, err
	}

	isAvailable, _ := p.CheckPackageAvailability(ctx, packageData, userMembership)
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
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] while GetPackageByID err : %v", err))

		return resp, err
	}

	if packageData.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = fmt.Sprintf("package with id %d not found", req.ID)
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] package with id %d not found", req.ID))

		return resp, fmt.Errorf("package with id %d not found", req.ID)
	}

	if req.DiscountedPrice != 0 && req.DiscountedPrice != packageData.DiscountedPrice {
		newDiscountExpired := getDiscountExpired(time.Now())
		err = p.PackageRepo.UpdatePackageDiscountExpired(ctx, newDiscountExpired)
		if err != nil {
			resp.Code = http.StatusBadRequest
			resp.Message = http.StatusText(http.StatusBadRequest)
			resp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdatePackage] while UpdatePackageDiscountExpired err : %v", err))

			return resp, err
		}
	}

	err = p.PackageRepo.UpdatePackage(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
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
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeletePackage] while GetPackageByID err : %v", err))

		return resp, err
	}

	if packageData.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
		resp.Error = fmt.Sprintf("package with id %d not found", packageId)
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeletePackage] package with id %d not found", packageId))

		return resp, fmt.Errorf("package with id %d not found", packageId)
	}

	err = p.PackageRepo.SoftDeletePackage(ctx, packageId, deletedBy)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = http.StatusText(http.StatusBadRequest)
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

func (p *PackageSvcImpl) CheckPackageAvailability(ctx context.Context, membershipPackage models.Package, userMembership models.UserMembership) (bool, error) {
	startDateStr, err := p.ConfigRepo.GetConfigByKey(utils.MembershipProgramStartDateConfig)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[CheckPackageAvailability] Failed to get %s: %v", utils.MembershipProgramStartDateConfig, err))
		return false, fmt.Errorf("failed to get %s: %w", utils.MembershipProgramStartDateConfig, err)
	}

	endDateStr, err := p.ConfigRepo.GetConfigByKey(utils.MembershipProgramEndDateConfig)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[CheckPackageAvailability] Failed to get %s: %v", utils.MembershipProgramEndDateConfig, err))
		return false, fmt.Errorf("failed to get %s: %w", utils.MembershipProgramEndDateConfig, err)
	}

	if startDateStr == "" || endDateStr == "" {
		slog.ErrorContext(ctx, "[CheckPackageAvailability] Start or End date config is empty")
		return true, nil
	}

	startDate, err := parseDate(startDateStr, time.DateOnly)
	if err != nil {
		return false, fmt.Errorf("failed to parse start date: %w", err)
	}

	endDate, err := parseDate(endDateStr, time.DateOnly)
	if err != nil {
		return false, fmt.Errorf("failed to parse end date: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	if today.Before(startDate) || today.After(endDate) {
		return true, nil
	}

	if userMembership.ID > 0 && isPreOrder(userMembership, startDate, endDate) {
		return checkPreOrderValidity(ctx, userMembership, membershipPackage, endDate)
	}

	return false, nil
}

func isPreOrder(um models.UserMembership, startDate, endDate time.Time) bool {
	switch um.Status {
	case models.MembershipStatusPreOrder:
		return true
	case models.MembershipStatusExpired:
		expDate, errParse := time.Parse(time.RFC3339, um.ExpiredAt)
		if errParse != nil {
			return false
		}
		return expDate.After(startDate) && expDate.Before(endDate)
	default:
		return false
	}
}

func parseDate(dateStr, layout string) (time.Time, error) {
	return time.Parse(layout, dateStr)

}

func checkPreOrderValidity(ctx context.Context, userMembership models.UserMembership, pkg models.Package, endDate time.Time) (bool, error) {
	expiredAt, err := time.Parse(time.RFC3339, userMembership.ExpiredAt)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[CheckPackageAvailability] Failed to parse expiredAt %s: %v", userMembership.ExpiredAt, err))
		return false, fmt.Errorf("failed to parse expiredAt: %w", err)
	}

	extendedExpiry := expiredAt.AddDate(0, int(pkg.DurationInMonth), 0)
	if extendedExpiry.After(endDate) {
		slog.ErrorContext(ctx, fmt.Sprintf("[CheckPackageAvailability] Extended expiry %s exceeds end date %s", extendedExpiry.Format(time.RFC3339), endDate.Format(time.RFC3339)))
		return false, nil
	}

	return true, nil
}
