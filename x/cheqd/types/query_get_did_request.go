package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (query *QueryGetDidRequest) Normalize() *QueryGetDidRequest {
	return &QueryGetDidRequest{
		Id: utils.NormalizeIdentifier(query.Id),
	}
}
