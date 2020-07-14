package src

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmwasm/cosmwasm-go/poc/std/ezdb"
)

func TestInit(t *testing.T){
	// we need to reset the database as this acts on a singleton
	ezdb.Reset()

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

func TestInitUsernameDifferentPassword(t *testing.T){
	// we need to reset the database as this acts on a singleton
	ezdb.Reset()

	msg := InitMsg{
		UserName: "AAAAAA",
		Password: "AAAAAA",
		Money:      1100,
	}

	res, err := go_init(msg)
	require.NotNil(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Err, `{"generic_err":{"msg":"error, UserName cannot equal with Password","backtrace":null}}`)
}

type QueryResponse struct {
	QueryResult string
}

func TestQuery(t *testing.T) {
	// we need to reset the database as this acts on a singleton
	ezdb.Reset()

	// initialize the contract
	msg := InitMsg{
		UserName: "AAAAAA",
		Password: "BBBBBB",
		Money:    1100,
	}
	_, err := go_init(msg)
	require.Nil(t, err)

	q := QueryMsg{
		QueryType: "balance",
	}
	res, err := go_query(q)
	require.Nil(t, err)
	require.NotNil(t, res)

	// output is json encoded base64
	var rawData []byte
	stderr := json.Unmarshal([]byte(res.Ok), &rawData)
	require.NoError(t, stderr)
	// this raw date is json-encoded QueryResponse
	var out QueryResponse
	stderr = json.Unmarshal(rawData, &out)
	require.NoError(t, stderr)

	// check the actual message
	require.Equal(t, out.QueryResult, "balance is : 1100")
}

func TestHandle(t *testing.T) {
	// we need to reset the database as this acts on a singleton
	ezdb.Reset()

	// initialize the contract
	msg := InitMsg{
		UserName:    "AAAAAA",
		Password: "BBBBBB",
		Money:      1100,
	}
	_, err := go_init(msg)
	require.Nil(t, err)

	// run a message
	h := HandleMsg{
		Operation: "save",
		Password:  msg.Password,
	}
	_, err = go_handle(h)
	require.Nil(t, err)

	// ensure balance has increased
	q := QueryMsg{
		QueryType: "balance",
	}
	res, err := go_query(q)
	require.Nil(t, err)
	require.NotNil(t, res)

	// output is json encoded base64
	var rawData []byte
	stderr := json.Unmarshal([]byte(res.Ok), &rawData)
	require.NoError(t, stderr)
	// this raw date is json-encoded QueryResponse
	var out QueryResponse
	stderr = json.Unmarshal(rawData, &out)
	require.NoError(t, stderr)

	// check the actual message
	require.Equal(t, out.QueryResult, "balance is : 1110")
}