package types

func NewService(id string, type_ string, serviceEndpoint string) *Service {
	return &Service{
		Id:              id,
		Type:            type_,
		ServiceEndpoint: serviceEndpoint,
	}
}
