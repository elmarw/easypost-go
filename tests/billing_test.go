package easypost_test

import "github.com/elmarw/easypost-go/v3"

func GetBillingMockRequests() []easypost.MockRequest {
	return []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/bank_accounts\\/\\S*\\/charges$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/credit_cards\\/\\S*\\/charges$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/bank_accounts\\/\\S*$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/credit_cards\\/\\S*$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/payment_methods$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "summary_123", "primary_payment_method": {"id": "card_123", "last4": "1234"}, "secondary_payment_method": {"id": "bank_123", "bank_name": "Mock Bank"}}`,
			},
		},
	}
}

func (c *ClientTests) TestDeletePaymentMethod() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	require := c.Require()

	err := client.DeletePaymentMethod(easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestFundWallet() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	require := c.Require()

	err := client.FundWallet("2000", easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestRetrievePaymentMethods() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	paymentMethods, err := client.RetrievePaymentMethods()
	require.NoError(err)

	assert.True(paymentMethods.PrimaryPaymentMethod != nil)
	assert.True(paymentMethods.SecondaryPaymentMethod != nil)
}
