package factory

import (
	"time"

	"hpc-express-service/auth"
	"hpc-express-service/common"
	"hpc-express-service/config"
	"hpc-express-service/customer"
	"hpc-express-service/dashboard"
	"hpc-express-service/gcs"
	inbound "hpc-express-service/inbound/express"
	"hpc-express-service/mawb"
	outboundExpress "hpc-express-service/outbound/express"
	outboundMawb "hpc-express-service/outbound/mawb"
	"hpc-express-service/setting"
	"hpc-express-service/ship2cu"
	"hpc-express-service/shopee"
	"hpc-express-service/tools/compare"
	"hpc-express-service/uploadlog"
	"hpc-express-service/user"
)

type ServiceFactory struct {
	AuthSvc                   auth.Service
	CommonSvc                 common.Service
	CompareSvc                compare.ExcelServiceInterface
	InboundExpressServiceSvc  inbound.InboundExpressService
	Ship2cuSvc                ship2cu.Service
	UploadlogSvc              uploadlog.Service
	OutboundExpressServiceSvc outboundExpress.OutboundExpressService
	OutboundMawbServiceSvc    outboundMawb.OutboundMawbService
	ShopeeSvc                 shopee.Service
	MawbSvc                   mawb.Service
	CustomerSvc               customer.Service
	DashboardSvc              dashboard.Service
	UserSvc                   user.Service
	SettingSvc                setting.Service
}

func NewServiceFactory(repo *RepositoryFactory, gcsClient *gcs.Client, conf *config.Config) *ServiceFactory {
	timeoutContext := time.Duration(60) * time.Second

	/*
	* Sharing Services
	 */

	// setting
	settingSvc := setting.NewService(
		repo.SettingRepo,
		timeoutContext,
	)

	// Ship2cu
	ship2cuSvc := ship2cu.NewService(
		repo.Ship2cuRepo,
		timeoutContext,
	)

	// Shopee
	shopeeSvc := shopee.NewService(
		repo.ShopeeRepo,
		timeoutContext,
	)

	// Upload Logging
	uploadlogSvc := uploadlog.NewService(
		repo.UploadlogRepo,
		timeoutContext,
		gcsClient,
	)

	// MAWB
	mawbSvc := mawb.NewService(
		repo.MawbRepo,
		timeoutContext,
	)

	// Customer
	customerSvc := customer.NewService(
		repo.CustomerRepo,
		timeoutContext,
	)

	// User
	userSvc := user.NewService(
		repo.UserRepo,
		timeoutContext,
	)
	/*
	* Sharing Services
	 */
	// Compare Service
	compareSvc := compare.NewExcelService(repo.CompareRepo)
	// Auth
	authSvc := auth.NewService(
		repo.AuthRepo,
		timeoutContext,
	)

	// Common
	dashboardSvc := dashboard.NewService(
		repo.DashboardRepo,
		timeoutContext,
	)

	// Common
	commonSvc := common.NewService(
		repo.CommonRepo,
		timeoutContext,
	)

	// Inbound Express
	inboundExpressServiceSvc := inbound.NewInboundExpressService(
		repo.InboundExpressRepositoryRepo,
		timeoutContext,
		ship2cuSvc,
		uploadlogSvc,
		repo.Ship2cuRepo,
	)

	// Outbound Express
	outboundExpressServiceSvc := outboundExpress.NewOutboundExpressService(
		repo.OutboundExpressRepositoryRepo,
		timeoutContext,
		shopeeSvc,
		uploadlogSvc,
	)

	// Outbound Mawb
	outboundMawbServiceSvc := outboundMawb.NewOutboundMawbService(
		repo.OutboundMawbRepositoryRepo,
		timeoutContext,
		gcsClient,
		conf,
	)

	return &ServiceFactory{
		AuthSvc:                   authSvc,
		CommonSvc:                 commonSvc,
		InboundExpressServiceSvc:  inboundExpressServiceSvc,
		Ship2cuSvc:                ship2cuSvc,
		UploadlogSvc:              uploadlogSvc,
		OutboundExpressServiceSvc: outboundExpressServiceSvc,
		OutboundMawbServiceSvc:    outboundMawbServiceSvc,
		ShopeeSvc:                 shopeeSvc,
		MawbSvc:                   mawbSvc,
		CustomerSvc:               customerSvc,
		DashboardSvc:              dashboardSvc,
		UserSvc:                   userSvc,
		CompareSvc:                compareSvc,
		SettingSvc:                settingSvc,
	}
}
