package factory

import (
	"time"

	"hpc-express-service/auth"
	"hpc-express-service/common"
	"hpc-express-service/customer"
	"hpc-express-service/dashboard"
	"hpc-express-service/gcs"
	inbound "hpc-express-service/inbound/express"
	"hpc-express-service/mawb"
	outbound "hpc-express-service/outbound/express"
	"hpc-express-service/shopee"
	"hpc-express-service/topgls"
	"hpc-express-service/uploadlog"
)

type ServiceFactory struct {
	AuthSvc                   auth.Service
	CommonSvc                 common.Service
	InboundExpressServiceSvc  inbound.InboundExpressService
	TopglsSvc                 topgls.Service
	UploadlogSvc              uploadlog.Service
	OutboundExpressServiceSvc outbound.OutboundExpressService
	ShopeeSvc                 shopee.Service
	MawbSvc                   mawb.Service
	CustomerSvc               customer.Service
	DashboardSvc              dashboard.Service
}

func NewServiceFactory(repo *RepositoryFactory, gcsClient *gcs.Client) *ServiceFactory {
	timeoutContext := time.Duration(60) * time.Second

	/*
	* Sharing Services
	 */

	// TOPGLS
	topglsSvc := topgls.NewService(
		repo.TopglsRepo,
		timeoutContext,
	)

	// TOPGLS
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
	/*
	* Sharing Services
	 */

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
		topglsSvc,
		uploadlogSvc,
	)

	// Outbound Express
	outboundExpressServiceSvc := outbound.NewOutboundExpressService(
		repo.OutboundExpressRepositoryRepo,
		timeoutContext,
		shopeeSvc,
		uploadlogSvc,
	)

	return &ServiceFactory{
		AuthSvc:                   authSvc,
		CommonSvc:                 commonSvc,
		InboundExpressServiceSvc:  inboundExpressServiceSvc,
		TopglsSvc:                 topglsSvc,
		UploadlogSvc:              uploadlogSvc,
		OutboundExpressServiceSvc: outboundExpressServiceSvc,
		ShopeeSvc:                 shopeeSvc,
		MawbSvc:                   mawbSvc,
		CustomerSvc:               customerSvc,
		DashboardSvc:              dashboardSvc,
	}
}
