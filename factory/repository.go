package factory

import (
	"hpc-express-service/auth"
	"hpc-express-service/common"
	"hpc-express-service/customer"
	"hpc-express-service/airline" // Added import
	"hpc-express-service/dashboard"
	inbound "hpc-express-service/inbound/express"
	"hpc-express-service/mawb"
	outbound "hpc-express-service/outbound/express"
	"hpc-express-service/shopee"
	"hpc-express-service/topgls"
	"hpc-express-service/uploadlog"
	"time"
)

type RepositoryFactory struct {
	AuthRepo                      auth.Repository
	CommonRepo                    common.Repository
	InboundExpressRepositoryRepo  inbound.InboundExpressRepository
	TopglsRepo                    topgls.Repository
	UploadlogRepo                 uploadlog.Repository
	OutboundExpressRepositoryRepo outbound.OutboundExpressRepository
	ShopeeRepo                    shopee.Repository
	MawbRepo                      mawb.Repository
	CustomerRepo                  customer.Repository
	DashboardRepo                 dashboard.Repository
	AirlineRepo                   airline.Repository // Added field
}

func NewRepositoryFactory() *RepositoryFactory {
	defaultTimeout := 5 * time.Second // Changed variable name and value

	return &RepositoryFactory{
		AuthRepo:                      auth.NewRepository(defaultTimeout), // Used defaultTimeout
		CommonRepo:                    common.NewRepository(defaultTimeout), // Used defaultTimeout
		InboundExpressRepositoryRepo:  inbound.NewInboundExpressRepository(defaultTimeout), // Used defaultTimeout
		TopglsRepo:                    topgls.NewRepository(defaultTimeout), // Used defaultTimeout
		UploadlogRepo:                 uploadlog.NewRepository(defaultTimeout), // Used defaultTimeout
		OutboundExpressRepositoryRepo: outbound.NewOutboundExpressRepository(defaultTimeout), // Used defaultTimeout
		ShopeeRepo:                    shopee.NewRepository(defaultTimeout), // Used defaultTimeout
		MawbRepo:                      mawb.NewRepository(defaultTimeout), // Used defaultTimeout
		CustomerRepo:                  customer.NewRepository(defaultTimeout), // Used defaultTimeout
		DashboardRepo:                 dashboard.NewRepository(defaultTimeout), // Used defaultTimeout
		AirlineRepo:                   airline.NewRepository(defaultTimeout),   // Added initialization
	}
}

// GetAirlineRepository returns the airline repository.
func (r *RepositoryFactory) GetAirlineRepository() airline.Repository {
	return r.AirlineRepo
}
