package dao

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dao := Dao{}
	err := dao.Init()
	require.Nil(t, err, "initialise es client should not have error")
}
