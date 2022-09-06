package utils_test

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func TestMustSplitDID(t *testing.T) {
	tests := []string{
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
		"did:cheqd:123456789abcdefg",
		"did:NOTcheqd:123456789abcdefg",
		"did:NOTcheqd:123456789abcdefg123456789abcdefg",
		"did:cheqd:testnet:123456789abcdefg",
	}

	_ = tests
}

func TestIsValidDID(t *testing.T) {
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

	convey.Convey("Given a set of DID variations", t, func() {
		for _, tt := range tests {
			convey.Convey(fmt.Sprintf("Given %s: %s", tt.name, tt.args.did), func() {
				isValid := utils.IsValidDID(tt.args.did, "cheqd", tt.args.allowedNamespaces)

				boolStr := "invalid"
				if isValid {
					boolStr = "valid"
				}

				convey.Convey(fmt.Sprintf("Then the DID should be %s", boolStr), func() {
					convey.So(isValid, convey.ShouldEqual, tt.expected)
				})
			})
		}
	})
}

func TestIsValidDIDGeneric(t *testing.T) {
	convey.Convey("Given a generic set of tests cases", t, func() {
		tests := []struct {
			text      string
			expected  bool
			did       string
			method    string
			allowedNS []string
		}{
			{"method and namespace are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet"}},
			{"method and namespaces are set", true, "did:cheqd:testnet:123456789abcdefg", "cheqd", []string{"testnet", "mainnet"}},
			{"method is not set", true, "did:cheqd:testnet:123456789abcdefg", "", []string{"testnet"}},
			{"method and namespaces are empty", true, "did:cheqd:testnet:123456789abcdefg", "", []string{}},
			{"namespace is absent in DID", true, "did:cheqd:123456789abcdefg", "", []string{}},
			{"method and namespaces are not set", false, "did:NOTcheqd:123456789abcdefg", "", []string{}},
		}

		for _, tt := range tests {
			convey.Convey(fmt.Sprintf("Given %s: %s", tt.text, tt.did), func() {
				isValid := utils.IsValidDID(tt.did, "cheqd", tt.allowedNS)

				boolStr := "invalid"
				if isValid {
					boolStr = "valid"
				}

				convey.Convey(fmt.Sprintf("Then the DID should be %s", boolStr), func() {
					convey.So(isValid, convey.ShouldEqual, tt.expected)
				})
			})
		}
	})
}
