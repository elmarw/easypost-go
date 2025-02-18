package easypost_test

import (
	"reflect"
	"strings"

	"github.com/elmarw/easypost-go/v3"
)

func (c *ClientTests) TestUserCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	userName := "Test User"

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test User", user.Name)
}

func (c *ClientTests) TestUserRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	authenticatedUser, err := client.RetrieveMe()
	require.NoError(err)

	retrievedUser, err := client.GetUser(authenticatedUser.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(retrievedUser))
	assert.True(strings.HasPrefix(retrievedUser.ID, "user_"))
}

func (c *ClientTests) TestUserRetrieveMe() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.RetrieveMe()
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
}

func (c *ClientTests) TestUserUpdate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.RetrieveMe()
	require.NoError(err)

	test_name := "Test User"

	updatedUser, err := client.UpdateUser(
		&easypost.UserOptions{
			ID:   user.ID,
			Name: &test_name,
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(updatedUser))
	assert.True(strings.HasPrefix(updatedUser.ID, "user_"))
	assert.Equal(test_name, updatedUser.Name)
}

func (c *ClientTests) TestUserUpdateBrand() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	color := "#123456"

	user, err := client.RetrieveMe()
	require.NoError(err)

	brand, err := client.UpdateBrand(
		map[string]interface{}{
			"color": color,
		},
		user.ID,
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Brand{}), reflect.TypeOf(brand))
	assert.True(strings.HasPrefix(brand.ID, "brd_"))
	assert.Equal(color, brand.Color)
}

func (c *ClientTests) TestUserDelete() {
	client := c.ProdClient()
	require := c.Require()

	userName := "Test User"

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	require.NoError(err)

	err = client.DeleteUser(user.ID)
	require.NoError(err)

	require.NoError(err)
}
