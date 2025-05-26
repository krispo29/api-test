package factory

import (
	"time"

	"hpc-express-service/airline" // Added import
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
	AirlineSvc                airline.Service // Added field
}

func NewServiceFactory(repoFactory *RepositoryFactory, gcsClient *gcs.Client) *ServiceFactory { // Renamed repo to repoFactory
	defaultTimeout := 5 * time.Second // Added defaultTimeout

	/*
	* Sharing Services
	 */

	// TOPGLS
	topglsSvc := topgls.NewService(
		repoFactory.TopglsRepo, // Used repoFactory
		defaultTimeout,         // Used defaultTimeout
	)

	// TOPGLS
	shopeeSvc := shopee.NewService(
		repoFactory.ShopeeRepo, // Used repoFactory
		defaultTimeout,         // Used defaultTimeout
	)

	// Upload Logging
	uploadlogSvc := uploadlog.NewService(
		repoFactory.UploadlogRepo, // Used repoFactory
		defaultTimeout,            // Used defaultTimeout
		gcsClient,
	)

	// MAWB
	mawbSvc := mawb.NewService(
		repoFactory.MawbRepo, // Used repoFactory
		defaultTimeout,       // Used defaultTimeout
	)

	// Customer
	customerSvc := customer.NewService(
		repoFactory.CustomerRepo, // Used repoFactory
		defaultTimeout,           // Used defaultTimeout
	)
	/*
	* Sharing Services
	 */

	// Auth
	authSvc := auth.NewService(
		repoFactory.AuthRepo, // Used repoFactory
		defaultTimeout,       // Used defaultTimeout
	)

	// Common
	dashboardSvc := dashboard.NewService(
		repoFactory.DashboardRepo, // Used repoFactory
		defaultTimeout,            // Used defaultTimeout
	)

	// Common
	commonSvc := common.NewService(
		repoFactory.CommonRepo, // Used repoFactory
		defaultTimeout,         // Used defaultTimeout
	)

	// Inbound Express
	inboundExpressServiceSvc := inbound.NewInboundExpressService(
		repoFactory.InboundExpressRepositoryRepo, // Used repoFactory
		defaultTimeout,                           // Used defaultTimeout
		topglsSvc,
		uploadlogSvc,
	)

	// Outbound Express
	outboundExpressServiceSvc := outbound.NewOutboundExpressService(
		repoFactory.OutboundExpressRepositoryRepo, // Used repoFactory
		defaultTimeout,                            // Used defaultTimeout
		shopeeSvc,
		uploadlogSvc,
	)

	// Airline
	airlineSvc := airline.NewService(
		repoFactory.GetAirlineRepository(), // Used repoFactory and new method
		defaultTimeout,                     // Used defaultTimeout
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
		AirlineSvc:                airlineSvc, // Added initialization
	}
}
