package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/logger"
	"github.com/metildachee/userie/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	setup()
	ctx := context.Background()
	_, err := NewDao(ctx)
	require.Nil(t, err, "initialise es client should not have error")
}

func TestCreateUser(t *testing.T) {
	setup()
	ctx := context.Background()
	u := models.User{
		Name:        "metchee",
		DOB:         int32(time.Now().Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")

	id, err := dao.Create(ctx, u)
	assert.Nil(t, err, "should not have error when indexing doc")
	assert.NotEqualValues(t, "0", id, "id should not be 0")
}

func TestCreateMultipleUsers(t *testing.T) {
	setup()
	ctx := context.Background()
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

	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")

	for i := 0; i < numOfUsersToCreate; i++ {
		wg.Add(1)
		go dao.Create(ctx, u, &wg)
	}
	id, err := dao.Create(ctx, u)
	require.Nil(t, err, "should not have error when create users")
	assert.NotEqualValues(t, "0", id, "id should not be 0")
	wg.Wait()

	users, err := dao.GetAll(ctx, 10)
	assert.Nil(t, err, "should not have error when get users")
	assert.True(t, len(users) >= numOfUsersToCreate+1, "we created 6 items, should have equal or more")
}

func TestBatchCreateUser(t *testing.T) {
	var (
		numOfUsers = 10
	)
	setup()
	ctx := context.Background()
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
	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")
	err = dao.BatchCreate(ctx, users)
	assert.Nil(t, err, "should not have error when create users")

	res, err := dao.GetAll(ctx, numOfUsers)
	assert.Nil(t, err, "should not have error when get users")
	assert.GreaterOrEqual(t, len(res), numOfUsers, "we created many items, should have equal or more")
}

func TestGetUser(t *testing.T) {
	setup()
	ctx := context.Background()
	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(ctx, "1")
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	fmt.Println("user", user)
}

func TestGetUsers(t *testing.T) {
	setup()
	var (
		minimumNumOfDocs = 5
	)
	ctx := context.Background()
	dao, err := NewDao(ctx)
	users, err := dao.GetAll(ctx, 10)
	assert.Nil(t, err, "should not have error when get users")
	assert.True(t, len(users) > minimumNumOfDocs)
}

func TestUpdateUser(t *testing.T) {
	setup()
	var (
		userId      = "1"
		updatedDesc = "let me change this up"
	)
	ctx := context.Background()
	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(ctx, userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	user.Description = updatedDesc
	err = dao.Update(ctx, user)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetById(ctx, user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotEqualValues(t, user.Description, updatedUser.Description, "should not have the same value")
}

func TestDeleteUser(t *testing.T) {
	setup()
	var (
		userId = "1"
	)
	ctx := context.Background()
	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(ctx, userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")

	err = dao.Delete(ctx, userId)
	assert.Nil(t, err, "should not have err when delete user")
	_, err = dao.GetById(ctx, userId)
	assert.NotNil(t, err)
}

func TestUpdateUserName(t *testing.T) {
	setup()
	var (
		userId      = "2"
		updatedName = "meow meow"
	)
	ctx := context.Background()
	dao, err := NewDao(ctx)
	assert.Nil(t, err, "should not have error when init")
	user, err := dao.GetById(ctx, userId)
	assert.Nil(t, err, "should not have err when getting user")
	assert.NotNil(t, user, "user should not be nil")
	err = dao.UpdateUserName(ctx, userId, updatedName)
	assert.Nil(t, err, "should not have err when update user")
	updatedUser, err := dao.GetById(ctx, user.ID)
	assert.Nil(t, err, "should not have err when getting user")
	assert.EqualValues(t, updatedName, updatedUser.Name, "should have the same value")
}

func setup() {
	lf, err := os.OpenFile("user_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("info logger", true, true, lf).Close()
}
