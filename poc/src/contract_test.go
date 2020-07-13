package src

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T){

	msg := InitMsg{
		UserName:    "AAAAAA",
		Password: "BBBBBB",
		Money:      1100,
	}

	res, err := go_init(msg)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.Ok, `{"messages":[],"log":[{"key":"result","value":"success to init contract"}],"data":null}`)
}

