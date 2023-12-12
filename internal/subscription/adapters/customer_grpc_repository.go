package adapters

import (
	"github.com/google/uuid"
	"gitlab.iman.uz/imandev/bnpl_payment/genproto/customer_service"
	"gitlab.iman.uz/imandev/bnpl_payment/internal/subscription/domain/customer"
	"golang.org/x/net/context"
)

type grpcCustomerRepository struct {
	grpcClient customer_service.CustomerServiceClient
}

func NewGrpcCustomerRepository(grpcClient customer_service.CustomerServiceClient) customer.Repository {
	return &grpcCustomerRepository{grpcClient: grpcClient}
}

func (s grpcCustomerRepository) FindOneByGuid(ctx context.Context, guid uuid.UUID) (*customer.Customer, error) {
	resp, err := s.grpcClient.GetCustomerDetails(ctx, &customer_service.GetCustomerDetailsRequest{
		CustomerGuid: guid.String(),
		Include:      "full_name",
	})

	if err != nil {
		return nil, err
	}

	customerId, err := uuid.Parse(resp.GetCustomerDetails().GetCustomer().GetGuid())
	if err != nil {
		return nil, err
	}

	return customer.NewCustomer(
		customerId,
		resp.CustomerDetails.Customer.FirstName,
		resp.CustomerDetails.Customer.LastName,
	), nil
}
