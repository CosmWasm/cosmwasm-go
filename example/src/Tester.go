package src

import (
	"bytes"
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
)

type Tester struct {
	deps *std.Extern
}

func (t Tester) DoTest() error {
	e := t.dbTest()
	if e != nil {
		return errors.New("Db test failed : " + e.Error())
	}
	return nil
}

func (t Tester) data() [][]byte {
	return [][]byte{
		[]byte("12345"), []byte("23456"), []byte("34567"), []byte("45678"),
	}
}

func (t Tester) dbTest() error {
	e := t.doWriteDbTest()
	if e != nil {
		return e
	}
	e = t.doReadDbTest()
	if e != nil {
		return e
	}
	e = t.doRemoveDbTest()
	if e != nil {
		return e
	}
	e = t.doRangeDbTest()
	if e != nil {
		return e
	}
	return nil
}

func (t Tester) doWriteDbTest() error {
	testKv := t.data()
	for _, kv := range testKv {
		e := t.deps.EStorage.Set(kv, kv)
		if e != nil {
			return errors.New("Testing write db : " + e.Error())
		}
	}
	return nil
}

func (t Tester) doReadDbTest() error {
	testKv := t.data()
	t.doWriteDbTest() //make sure we have right test data in db
	for _, kv := range testKv {
		v, e := t.deps.EStorage.Get(kv)
		if e != nil {
			return errors.New("Testing write db : " + e.Error())
		}
		if !bytes.Equal(v, kv) {
			return errors.New("Testing Read db, got wrong data, excepted : " + string(kv) + ",got " + string(v))
		}
	}
	return nil
}

func (t Tester) doRemoveDbTest() error {
	testKv := t.data()
	for _, kv := range testKv {
		e := t.deps.EStorage.Remove(kv)
		if e != nil {
			return errors.New("Testing write db : " + e.Error())
		}
		v, e := t.deps.EStorage.Get(kv)
		//e must not nil
		if e == nil {
			return errors.New("Deleted values has read out : " + string(v))
		}
	}

	return nil
}

func (t Tester) doRangeDbTest() error {
	t.doWriteDbTest() //make sure we have right test data in db
	return nil
}
