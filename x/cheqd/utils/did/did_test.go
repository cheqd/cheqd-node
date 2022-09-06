package did_test

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func TestIsValidDIDSpec(t *testing.T) {
	t.Parallel()

	type args struct {
		did               string
		allowedNamespaces []string
	}

	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "indy-style did (32)",
			args: args{
				did:               "did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
				allowedNamespaces: []string{"mainnet"},
			},
			expected: true,
		},
		{
			name: "indy-style did (16)",
			args: args{
				did:               "did:cheqd:testnet:DAzMQo4MDMxCjgwM",
				allowedNamespaces: []string{"testnet"},
			},
			expected: true,
		},
		{
			name: "indy-style did (no namespace)",
			args: args{
				did:               "did:cheqd:6cgbu8ZPoWTnR5Rv",
				allowedNamespaces: []string{},
			},
			expected: true,
		},
		{
			name: "uuid-style did",
			args: args{
				did:               "did:cheqd:mainnet:de9786cd-ec53-458c-857c-9342cf264f80",
				allowedNamespaces: []string{"mainnet"},
			},
			expected: true,
		},
		{
			name: "uuid-style did (no namespace)",
			args: args{
				did:               "did:cheqd:de9786cd-ec53-458c-857c-9342cf264f80",
				allowedNamespaces: []string{},
			},
			expected: true,
		},
		{
			name: "did with key",
			args: args{
				did:               "did:cheqd:testnet:DAzMQo4MDMxCjgwM",
				allowedNamespaces: []string{},
			},
			expected: true,
		},
		{
			name: "did with namespace not allowed",
			args: args{
				did:               "did:cheqd:testnet:DAzMQo4MDMxCjgwM",
				allowedNamespaces: []string{"mainnet"},
			},
			expected: false,
		},
		{
			name: "did with invalid method",
			args: args{
				did:               "did:notcheqd:testnet:DAzMQo4MDMxCjgwM",
				allowedNamespaces: []string{"testnet"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		convey.Convey(fmt.Sprintf("Given %s: %s", tt.name, tt.args.did), t, func() {
			isValid := utils.IsValidDID(tt.args.did, "cheqd", tt.args.allowedNamespaces)

			var not string
			if !isValid {
				not = "not"
			}

			convey.Convey(fmt.Sprintf("Then the DID should %s be valid", not), func() {
				if tt.expected {
					convey.So(isValid, convey.ShouldBeTrue)
				} else {
					convey.So(isValid, convey.ShouldBeFalse)
				}
			})
		})
	}
}
