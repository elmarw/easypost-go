package easypost

import (
	"context"
	"net/http"
)

// AddressVerificationFieldError provides additional information on address
// validation errors.
type AddressVerificationFieldError struct {
	Code       string `json:"code,omitempty"`
	Field      string `json:"field,omitempty"`
	Message    string `json:"message,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// AddressVerificationDetails contains extra information related to address
// verification.
type AddressVerificationDetails struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeZone  string  `json:"time_zone"`
}

// AddressVerification holds data relating to address verification status.
type AddressVerification struct {
	Success bool                             `json:"success"`
	Errors  []*AddressVerificationFieldError `json:"errors"`
	Details *AddressVerificationDetails      `json:"details"`
}

// AddressVerifications contains the result of the requested address
// verifications.
type AddressVerifications struct {
	ZIP4     *AddressVerification `json:"zip4"`
	Delivery *AddressVerification `json:"delivery"`
}

// Address objects are used to represent people, places, and organizations in a
// number of contexts.
type Address struct {
	ID              string                `json:"id,omitempty"`
	Object          string                `json:"object,omitempty"`
	Reference       string                `json:"reference,omitempty"`
	Mode            string                `json:"mode,omitempty"`
	CreatedAt       *DateTime             `json:"created_at,omitempty"`
	UpdatedAt       *DateTime             `json:"updated_at,omitempty"`
	Street1         string                `json:"street1,omitempty"`
	Street2         string                `json:"street2,omitempty"`
	City            string                `json:"city,omitempty"`
	State           string                `json:"state,omitempty"`
	Zip             string                `json:"zip,omitempty"`
	Country         string                `json:"country,omitempty"`
	Name            string                `json:"name,omitempty"`
	Company         string                `json:"company,omitempty"`
	Phone           string                `json:"phone,omitempty"`
	Email           string                `json:"email,omitempty"`
	Residential     bool                  `json:"residential,omitempty"`
	CarrierFacility string                `json:"carrier_facility,omitempty"`
	FederalTaxID    string                `json:"federal_tax_id,omitempty"`
	StateTaxID      string                `json:"state_tax_id,omitempty"`
	Verifications   *AddressVerifications `json:"verifications,omitempty"`
}

// CreateAddressOptions is used to specify verification options for address
// creation.
type CreateAddressOptions struct {
	Verify       []string `json:"verify,omitempty"`
	VerifyStrict []string `json:"verify_strict,omitempty"`
}

type createAddressRequest struct {
	*CreateAddressOptions
	Address *Address `json:"address,omitempty"`
}

type AddressVerifyResponse struct {
	Address *Address `json:"address,omitempty"`
}

// ListAddressResult holds the results from the list addresses API.
type ListAddressResult struct {
	Addresses []*Address `json:"addresses,omitempty"`
	PaginatedCollection
}

// For some reason, the verify API returns the address in a nested dictionary.
type verifyAddressResponse struct {
	Address **Address `json:"address,omitempty"`
}

// CreateAddress submits a request to create a new address, and returns the
// result.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateAddress(
//		&easypost.Address{
//			Street1: "417 Montgomery Street",
//			Street2: "Floor 5",
//			City:    "San Francisco",
//			State:   "CA",
//			Zip:     "94104",
//			Country: "US",
//			Company: "EasyPost",
//			Phone:   "415-123-4567",
//		},
//		&CreateAddressOptions{Verify: []string{"delivery"}},
//	)
func (c *Client) CreateAddress(in *Address, opts *CreateAddressOptions) (out *Address, err error) {
	return c.CreateAddressWithContext(context.Background(), in, opts)
}

// CreateAddressWithContext performs the same operation as CreateAddress, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateAddressWithContext(ctx context.Context, in *Address, opts *CreateAddressOptions) (out *Address, err error) {
	req := &createAddressRequest{CreateAddressOptions: opts, Address: in}
	err = c.post(ctx, "addresses", req, &out)
	return
}

// ListAddresses provides a paginated result of Address objects.
func (c *Client) ListAddresses(opts *ListOptions) (out *ListAddressResult, err error) {
	return c.ListAddressesWithContext(context.Background(), opts)
}

// ListAddressesWithContext performs the same operation as ListAddresses, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListAddressesWithContext(ctx context.Context, opts *ListOptions) (out *ListAddressResult, err error) {
	err = c.do(ctx, http.MethodGet, "addresses", c.convertOptsToURLValues(opts), &out)
	return
}

// GetNextAddressPage returns the next page of addresses
func (c *Client) GetNextAddressPage(collection *ListAddressResult) (out *ListAddressResult, err error) {
	return c.GetNextAddressPageWithContext(context.Background(), collection)
}

// GetNextAddressPageWithPageSize returns the next page of addresses with a specific page size
func (c *Client) GetNextAddressPageWithPageSize(collection *ListAddressResult, pageSize int) (out *ListAddressResult, err error) {
	return c.GetNextAddressPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextAddressPageWithContext performs the same operation as GetNextAddressPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextAddressPageWithContext(ctx context.Context, collection *ListAddressResult) (out *ListAddressResult, err error) {
	return c.GetNextAddressPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextAddressPageWithPageSizeWithContext performs the same operation as GetNextAddressPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextAddressPageWithPageSizeWithContext(ctx context.Context, collection *ListAddressResult, pageSize int) (out *ListAddressResult, err error) {
	if len(collection.Addresses) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.Addresses[len(collection.Addresses)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListAddressesWithContext(ctx, params)
}

// VerifyAddress performs address verification.
func (c *Client) VerifyAddress(addressID string) (out *Address, err error) {
	return c.VerifyAddressWithContext(context.Background(), addressID)
}

// VerifyAddressWithContext performs the same operation as VerifyAddress, but
// allows specifying a context that can interrupt the request.
func (c *Client) VerifyAddressWithContext(ctx context.Context, addressID string) (out *Address, err error) {
	path := "addresses/" + addressID + "/verify"
	res := &verifyAddressResponse{Address: &out}
	err = c.get(ctx, path, res)
	return
}

// GetAddress retrieves a previously-created address by its ID.
func (c *Client) GetAddress(addressID string) (out *Address, err error) {
	return c.GetAddressWithContext(context.Background(), addressID)
}

// GetAddressWithContext performs the same operation as GetAddress, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetAddressWithContext(ctx context.Context, addressID string) (out *Address, err error) {
	err = c.get(ctx, "addresses/"+addressID, &out)
	return
}

// CreateAndVerifyAddress Create Address object and immediately verify it.
func (c *Client) CreateAndVerifyAddress(in *Address, opts *CreateAddressOptions) (out *Address, err error) {
	return c.CreateAndVerifyAddressWithContext(context.Background(), in, opts)
}

// CreateAndVerifyAddressWithContext performs the same operation as CreateAndVerifyAddress, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateAndVerifyAddressWithContext(ctx context.Context, in *Address, opts *CreateAddressOptions) (out *Address, err error) {
	req := &createAddressRequest{CreateAddressOptions: opts, Address: in}
	response := AddressVerifyResponse{}
	err = c.post(ctx, "addresses/create_and_verify", req, &response)
	out = response.Address
	return
}
