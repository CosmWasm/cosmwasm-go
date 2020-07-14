package src

import (
	"strconv"

	"github.com/cosmwasm/cosmwasm-go/poc/std"
	"github.com/cosmwasm/cosmwasm-go/poc/std/ezdb"
)

type InitMsg struct {
	UserName string
	Password string
	Money    int
}

type HandleMsg struct {
	Operation string
	Password  string
}

type QueryMsg struct {
	QueryType string
}

type QueryResponse struct {
	QueryResult string
}

func getMoneyLeft() (int, error) {
	money, err := ezdb.ReadStorage([]byte("Money"))
	if err != nil {
		return 0, err
	}
	moneyInt, e := strconv.Atoi(string(money))
	if e != nil {
		return 0, e
	}
	return moneyInt, nil
}

func saveMoeny(money int) error {
	return ezdb.WriteStorage([]byte("Money"), []byte(strconv.Itoa(money)))
}

func go_init(msg InitMsg) (*OKResponse, *ERRResponse) {
	if msg.UserName == msg.Password {
		return nil, newErrResponse(std.Build_ErrResponse("error, UserName cannot equal with Password"))
	}
	_, err := ezdb.ReadStorage([]byte("inited"))
	if err == nil {
		return nil, newErrResponse(std.Build_ErrResponse("Contract already inited"))
	}
	e := ezdb.WriteStorage([]byte("UserName"), []byte(msg.UserName))
	if e != nil {
		return nil, newErrResponse(std.Build_ErrResponse("WriteStorage"))
	}
	ezdb.WriteStorage([]byte("Password"), []byte(msg.Password))
	ezdb.WriteStorage([]byte("Money"), []byte(strconv.Itoa(msg.Money)))
	return newOkResponse(std.Build_OkResponse("success to init contract")), nil
}

func go_handle(msg HandleMsg) (*OKResponse, *ERRResponse) {
	p, err := ezdb.ReadStorage([]byte("Password"))
	if err != nil {
		return nil, newErrResponse(std.Build_ErrResponse(err.Error()))
	}
	if msg.Password != string(p) {
		return nil, newErrResponse(std.Build_ErrResponse("Wrong password, check again"))
	}
	switch msg.Operation {
	case "burn":
		moneyInt, e := getMoneyLeft()
		if e != nil {
			return nil, newErrResponse(std.Build_ErrResponse(e.Error()))
		}
		if moneyInt == 0 {
			return nil, newErrResponse(std.Build_ErrResponse("Sorry, all money has burned~ try to save some money"))
		}
		if moneyInt < 10 && moneyInt > 0 {
			moneyInt = 0
		} else {
			moneyInt -= 10
		}
		e = saveMoeny(moneyInt)
		if e != nil {
			return nil, newErrResponse(std.Build_ErrResponse(err.Error()))
		}
	case "save":
		moneyInt, e := getMoneyLeft()
		if e != nil {
			return nil, newErrResponse(std.Build_ErrResponse(err.Error()))
		}
		moneyInt += 10
		e = saveMoeny(moneyInt)
		if e != nil {
			return nil, newErrResponse(std.Build_ErrResponse(err.Error()))
		}
	default:
		return nil, newErrResponse(std.Build_ErrResponse("Unsupport operation :" + msg.Operation))
	}
	return newOkResponse(std.Build_OkResponse("handle run success")), nil
}

func go_query(msg QueryMsg) (*OKResponse, *ERRResponse) {
	switch msg.QueryType {
	case "balance":
		moneyInt, e := getMoneyLeft()
		if e != nil {
			return nil, newErrResponse(std.Build_ErrResponse(e.Error()))
		}
		resp := QueryResponse{"balance is : " + strconv.Itoa(moneyInt)}
		return EncodeResult(resp)
	case "user":
		username, err := ezdb.ReadStorage([]byte("UserName"))
		if err != nil {
			return nil, newErrResponse(std.Build_ErrResponse("Read UserName failed: " + err.Error()))
		}
		resp := QueryResponse{string(username)}
		return EncodeResult(resp)
	default:
		return nil, newErrResponse(std.Build_ErrResponse("required balance or user, found unsupport query type :" + msg.QueryType))
	}
}
