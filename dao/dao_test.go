package dao

import (
	"testing"
	"time"

	"github.com/metildachee/userie/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dao := Dao{}
	err := dao.Init()
	require.Nil(t, err, "initialise es client should not have error")
}

func TestCreateUser(t *testing.T) {
	u := model.User{
		ID:          1,
		Name:        "metchee",
		DOB:         int32(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	dao := NewDao()
	err := dao.Init()
	assert.Nil(t, err, "should not have error when init")

	err = dao.CreateUser(u)
	assert.Nil(t, err, "should not have error when indexing doc")
}
