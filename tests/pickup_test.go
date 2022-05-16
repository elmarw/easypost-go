package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestPickupCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(pickup))
	assert.True(strings.HasPrefix(pickup.ID, "pickup_"))
	assert.NotNil(pickup.PickupRates)
}

func (c *ClientTests) TestPickupRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievePickup, err := client.GetPickup(pickup.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(retrievePickup))
	assert.Equal(pickup, retrievePickup)
}

func (c *ClientTests) TestPickupBuy() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	if err != nil {
		c.T().Error(err)
		return
	}

	boughtPickup, err := client.BuyPickup(
		pickup.ID,
		&easypost.PickupRate{
			Carrier: c.fixture.USPS(),
			Service: c.fixture.PickupService(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(boughtPickup))
	assert.True(strings.HasPrefix(boughtPickup.ID, "pickup_"))
	assert.NotNil(pickup.Confirmation)
	assert.Equal("scheduled", boughtPickup.Status)
}

func (c *ClientTests) TestPickupCancel() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	if err != nil {
		c.T().Error(err)
		return
	}

	boughtPickup, err := client.BuyPickup(
		pickup.ID,
		&easypost.PickupRate{
			Carrier: c.fixture.USPS(),
			Service: c.fixture.PickupService(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	cancelledPickup, err := client.CancelPickup(boughtPickup.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(cancelledPickup))
	assert.True(strings.HasPrefix(cancelledPickup.ID, "pickup_"))
	assert.Equal("canceled", cancelledPickup.Status)
}
