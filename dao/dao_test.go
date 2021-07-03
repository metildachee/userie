package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/metildachee/userie/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	_, err := NewDao()
	require.Nil(t, err, "initialise es client should not have error")
}

func TestCreateUser(t *testing.T) {
	u := model.User{
		Name:        "metchee",
		DOB:         float64(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       float64(time.Now().Unix()),
	}

	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")

	err = dao.CreateUser(u)
	assert.Nil(t, err, "should not have error when indexing doc")
}

func TestGetUser(t *testing.T) {
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetUser("4H59bHoBOTxfBxWqWN9B")
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	fmt.Println("user", user)
}