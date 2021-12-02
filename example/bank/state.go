package bank

import (
	"errors"
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/storage"
)

var (
	ErrDenomMetadataAlreadyExists = errors.New("bank: DenomMetadata already exists")
	ErrDenomMetadataNotFound      = errors.New("bank: DenomMetadata not found")
)

var DenomMetadataState = DenomMetadataStateHandler{ns: storage.NewNamespaced("denom_metadata")}

// DenomMetadata defines the denom metadata
//go:generate ../../bin/generator state map DenomMetadata
type DenomMetadata struct {
	Name    string `json:"name" orm:"primaryKey,1"`
	Decimal int
}

func (d DenomMetadata) MarshalJSON() ([]byte, error) {
	panic("impl")
}

func (d *DenomMetadata) UnmarshalJSON(b []byte) error {
	panic("impl")
}

type DenomMetadataStateHandler struct {
	ns storage.Namespaced
}

// Create creates a new DenomMetadata. Returns ErrDenomMetadataAlreadyExists if the object is already found.
func (x DenomMetadataStateHandler) Create(s std.Storage, denomMetadata DenomMetadata) error {
	key := []byte(denomMetadata.Name)
	v := s.Get(key)
	if v != nil {
		return ErrDenomMetadataAlreadyExists
	}

	v, err := denomMetadata.MarshalJSON()
	if err != nil {
		return err
	}

	s.Set(key, v)
	return nil
}

func (x DenomMetadataStateHandler) Read(s std.Storage, name string) (DenomMetadata, error) {
	key := []byte(name)
	v := s.Get(key)
	if v == nil {
		return DenomMetadata{}, ErrDenomMetadataNotFound
	}
	denomMetadata := DenomMetadata{}
	err := denomMetadata.UnmarshalJSON(v)
	if err != nil {
		return DenomMetadata{}, err
	}

	return denomMetadata, nil
}

func (x DenomMetadataStateHandler) Update(s std.Storage, name string, f func(old *DenomMetadata) (updated *DenomMetadata, err error)) (err error) {
	key := []byte(name)
	v := s.Get(key)
	var old *DenomMetadata
	// unmarshal state only if we find something
	if v != nil {
		old = new(DenomMetadata)
		err = old.UnmarshalJSON(v)
		if err != nil {
			return err
		}
	}

	updated, err := f(old)
	if err != nil {
		return err
	}

	// if no updates just ignore
	if updated == nil {
		return nil
	}
	// otherwise marshal and save
	v, err = updated.MarshalJSON()
	if err != nil {
		return err
	}

	s.Set(key, v)

	return nil
}

func (x DenomMetadataStateHandler) Delete(s std.Storage, denomMetadata DenomMetadata) error {
	key := []byte(denomMetadata.Name)
	v := s.Get(key)
	if v == nil {
		return ErrDenomMetadataNotFound
	}

	s.Remove(key)
	return nil
}
