package elasticsearch

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/metildachee/userie/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	_, err := NewDao()
	require.Nil(t, err, "initialise es client should not have error")
}

func TestCreateUser(t *testing.T) {
	u := models.User{
		Name:        "metchee",
		DOB:         int32(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")

	id, err := dao.Create(u)
	assert.Nil(t, err, "should not have error when indexing doc")
	assert.NotEqualValues(t, "0", id, "id should not be 0")
}

func TestCreateMultipleUsers(t *testing.T) {
	var (
		numOfUsersToCreate = 5
		wg                 sync.WaitGroup
	)
	u := models.User{
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
		go dao.Create(u, &wg)
	}
	id, err := dao.Create(u)
	require.Nil(t, err, "should not have error when create users")
	assert.NotEqualValues(t, "0", id, "id should not be 0")
	wg.Wait()

	users, err := dao.GetAll(10)
	assert.Nil(t, err, "should not have error when get users")
	assert.True(t, len(users) >= numOfUsersToCreate+1, "we created 6 items, should have equal or more")
}

func TestBatchCreateUser(t *testing.T) {
	var (
		numOfUsers = 10
	)
	users := make([]models.User, 0)
	for i := 0; i < numOfUsers; i++ {
		u := models.User{
			Name:        fmt.Sprintf("metchee %d", i),
			DOB:         int32(time.Now().Unix()),
			Address:     fmt.Sprintf("kent ridge %d", i),
			Description: fmt.Sprintf("default user info %d", i),
			Ctime:       int32(time.Now().Unix()),
		}
		users = append(users, u)
	}
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	err = dao.BatchCreate(users)
	assert.Nil(t, err, "should not have error when create users")

	res, err := dao.GetAll(numOfUsers)
	assert.Nil(t, err, "should not have error when get users")
	assert.GreaterOrEqual(t, len(res), numOfUsers, "we created many items, should have equal or more")
}

func TestGetUser(t *testing.T) {
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById("1")
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	fmt.Println("user", user)
}

func TestGetUsers(t *testing.T) {
	var (
		minimumNumOfDocs = 5
	)
	dao, err := NewDao()
	users, err := dao.GetAll(10)
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
	user, err := dao.GetById(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	user.Description = updatedDesc
	err = dao.Update(user)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetById(user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotEqualValues(t, user.Description, updatedUser.Description, "should not have the same value")
}

func TestDeleteUser(t *testing.T) {
	var (
		userId = "1"
	)
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	err = dao.Delete(userId)
	assert.Nil(t, err, "should not have err when delete user")
	_, err = dao.GetById(userId)
	assert.NotNil(t, err)
}

func TestUpdateUserName(t *testing.T) {
	var (
		userId      = "2"
		updatedName = "meow meow"
	)
	dao, err := NewDao()
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	err = dao.UpdateUserName(userId, updatedName)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetById(user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.EqualValues(t, updatedName, updatedUser.Name, "should not have the same value")
}
