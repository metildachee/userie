package dao

import (
	"fmt"
	"sync"
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
		DOB:         int32(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")

	err = dao.CreateUser(u)
	assert.Nil(t, err, "should not have error when indexing doc")
}

func TestCreateMultipleUsers(t *testing.T) {
	var (
		numOfUsersToCreate = 5
		wg                 sync.WaitGroup
	)
	u := model.User{
		Name:        "metchee",
		DOB:         int32(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")

	for i := 0; i < numOfUsersToCreate; i++ {
		wg.Add(1)
		go dao.CreateUser(u, &wg)
	}
	err = dao.CreateUser(u)
	assert.Nil(t, err, "should not have error when create users")
	wg.Wait()

	users, err := dao.GetUsers(10)
	assert.Nil(t, err, "should not have error when get users")
	assert.True(t, len(users) >= numOfUsersToCreate+1, "we created 6 items, should have equal or more")
}

func TestGetUser(t *testing.T) {
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetUser("1")
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	fmt.Println("user", user)
}

func TestGetUsers(t *testing.T) {
	var (
		minimumNumOfDocs = 5
	)
	dao, err := NewDao()
	users, err := dao.GetUsers(10)
	assert.Nil(t, err, "should not have error when get users")
	assert.True(t, len(users) > minimumNumOfDocs)
}

func TestUpdateUser(t *testing.T) {
	var (
		userId      = "1"
		updatedDesc = "let me change this up"
	)
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetUser(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	user.Description = updatedDesc
	err = dao.UpdateUser(user)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetUser(user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotEqualValues(t, user.Description, updatedUser.Description, "should not have the same value")
}

func TestDeleteUser(t *testing.T) {
	var (
		userId = "1"
	)
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetUser(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	err = dao.DeleteUser(userId)
	assert.Nil(t, err, "should not have err when delete user")
	_, err = dao.GetUser(userId)
	assert.NotNil(t, err)
}

func TestUpdateUserName(t *testing.T) {
	var (
		userId      = "1"
		updatedName = "meow meow"
	)
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetUser(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	time.Sleep(5 * time.Second)

	err = dao.UpdateUserName(userId, updatedName)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetUser(user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.EqualValues(t, updatedUser.Name, updatedName, "should not have the same value")
}
