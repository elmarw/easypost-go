package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestShipmentCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.NotNil(shipment.Rates)
	assert.Equal("PNG", shipment.Options.LabelFormat)
	assert.Equal("123", shipment.Options.InvoiceNumber)
	assert.Equal("123", shipment.Reference)
}

func (c *ClientTests) TestShipmentRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	retrievedShipment, err := client.GetShipment(shipment.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(retrievedShipment))
	assert.Equal(shipment, retrievedShipment)
}

func (c *ClientTests) TestShipmentAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipments, err := client.ListShipments(
		&easypost.ListShipmentsOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	shipmentsList := shipments.Shipments

	assert.LessOrEqual(len(shipmentsList), c.fixture.pageSize())
	assert.NotNil(shipments.HasMore)
	for _, shipment := range shipmentsList {
		assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	}
}

func (c *ClientTests) TestShipmentBuy() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	lowestRate, err := client.LowestRate(shipment)
	require.NoError(err)

	boughtShipment, err := client.BuyShipment(shipment.ID, &lowestRate, "")
	require.NoError(err)

	assert.NotNil(boughtShipment.PostageLabel)
}

func (c *ClientTests) TestShipmentRegenerateRates() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	rates, err := client.RerateShipment(shipment.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf([]*easypost.Rate{}), reflect.TypeOf(rates))
	for _, rate := range rates {
		assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	}
}

func (c *ClientTests) TestShipmentConvertLabel() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	shipmentWithNewLabel, err := client.GetShipmentLabel(shipment.ID, "ZPL")
	require.NoError(err)

	assert.NotNil(shipmentWithNewLabel.PostageLabel.LabelZPLURL)
}

// If the shipment was purchased with a USPS rate, it must have had its insurance set to `0` when bought
// so that USPS doesn't automatically insure it so we could manually insure it here.
func (c *ClientTests) TestShipmentInsure() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.OneCallBuyShipment()
	// Set to 0 so USPS doesn't insure this automatically and we can insure the shipment manually
	shipmentData.Insurance = "0"

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	insuredShipment, err := client.InsureShipment(shipment.ID, "100")
	require.NoError(err)

	assert.Equal("100.00", insuredShipment.Insurance)
}

// Refunding a test shipment must happen within seconds of the shipment being created as test shipments naturally
// follow a flow of created -> delivered to cycle through tracking events in test mode - as such anything older
// than a few seconds in test mode may not be refundable.
func (c *ClientTests) TestShipmentRefund() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	refundShipment, err := client.RefundShipment(shipment.ID)
	require.NoError(err)

	assert.Equal("submitted", refundShipment.RefundStatus)
}

func (c *ClientTests) TestShipmentSmartrate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	assert.NotNil(shipment.Rates)

	smartrates, err := client.GetShipmentSmartrates(shipment.ID)
	require.NoError(err)

	assert.Equal(shipment.Rates[0].ID, smartrates[0].ID)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile50)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile75)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile85)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile90)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile95)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile97)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile99)
}

func (c *ClientTests) TestShipmentCreateEmptyObjects() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.BasicShipment()
	shipmentData.Options = nil
	shipmentData.TaxIdentifiers = nil
	shipmentData.Reference = ""
	shipmentData.CustomsInfo = &easypost.CustomsInfo{}
	shipmentData.CustomsInfo.CustomsItems = []*easypost.CustomsItem{}

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.NotNil(shipment.Options) // The EasyPost API populates some default values here
	assert.Nil(shipment.CustomsInfo)
	assert.Nil(shipment.TaxIdentifiers)
	assert.NotNil(shipment.Reference)
}

func (c *ClientTests) TestShipmentCreateTaxIdentifier() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.BasicShipment()
	shipmentData.TaxIdentifiers = []*easypost.TaxIdentifier{c.fixture.TaxIdentifier()}

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.Equal("IOSS", shipment.TaxIdentifiers[0].TaxIdType)
}

func (c *ClientTests) TestShipmentCreateWithIds() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	fromAddress, err := client.CreateAddress(c.fixture.BasicAddress(), nil)
	require.NoError(err)
	toAddress, err := client.CreateAddress(c.fixture.BasicAddress(), nil)
	require.NoError(err)
	parcel, err := client.CreateParcel(c.fixture.BasicParcel())
	require.NoError(err)

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			FromAddress: &easypost.Address{ID: fromAddress.ID},
			ToAddress:   &easypost.Address{ID: toAddress.ID},
			Parcel:      &easypost.Parcel{ID: parcel.ID},
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.True(strings.HasPrefix(shipment.FromAddress.ID, "adr_"))
	assert.True(strings.HasPrefix(shipment.ToAddress.ID, "adr_"))
	assert.True(strings.HasPrefix(shipment.Parcel.ID, "prcl_"))
	assert.Equal("388 Townsend St", shipment.FromAddress.Street1)
}
