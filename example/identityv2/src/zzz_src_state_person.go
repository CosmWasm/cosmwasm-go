package src

import (
	errors "errors"
	std "github.com/cosmwasm/cosmwasm-go/std"
	storage "github.com/cosmwasm/cosmwasm-go/std/storage"
	keys "github.com/cosmwasm/cosmwasm-go/std/storage/keys"
)

var (
	// ErrPersonNotFound is returned when the object is not found.
	ErrPersonNotFound = errors.New("src: Person not found")
	// ErrPersonAlreadyExists is returned when the object already exists.
	ErrPersonAlreadyExists = errors.New("src: Person already exists")
)

var (
	// PersonState provides access to the Person namespaced storage.
	PersonState = PersonStateHandler{ns: storage.NewNamespaced("person")}
)

// PersonStateHandler provides ORM functionality around Person.
type PersonStateHandler struct {
	ns storage.Namespaced
}

// Create handles creation of Person objects.
// Returns ErrPersonAlreadyExists in case the object exists.
func (x PersonStateHandler) Create(storage std.Storage, o Person) error {
	_k := keys.StringPrimaryKey(o.Address)
	_v := x.ns.Get(storage, _k)
	if _v != nil {
		return ErrPersonAlreadyExists
	}
	_b, err := o.MarshalJSON()
	if err != nil {
		return err
	}
	x.ns.Set(storage, _k, _b)
	return nil
}

// Read returns Person given its Address.
// Returns ErrPersonNotFound in case nothing is found.
func (x PersonStateHandler) Read(storage std.Storage, string string) (Person, error) {
	_k := keys.StringPrimaryKey(string)
	_v := x.ns.Get(storage, _k)
	if _v == nil {
		return Person{}, ErrPersonNotFound
	}
	_o := new(Person)
	err := _o.UnmarshalJSON(_v)
	if err != nil {
		return Person{}, err
	}
	return *_o, nil
}

// Update updates an instance of Person, given its Address by running the provided function f.
// If f passes a nil Person it means the object was not found.
// Returning a nil Person from f will cause no updates in the object.
func (x PersonStateHandler) Update(storage std.Storage, string string, f func(old *Person) (*Person, error)) (err error) {
	_k := keys.StringPrimaryKey(string)
	_v := x.ns.Get(storage, _k)

	var _old *Person

	if _v != nil {
		_old = new(Person)
		err = _old.UnmarshalJSON(_v)
		if err != nil {
			return err
		}
	}

	_updated, err := f(_old)
	if err != nil {
		return err
	}

	if _updated == nil {
		return nil
	}

	_v, err = _updated.MarshalJSON()
	if err != nil {
		return err
	}

	x.ns.Set(storage, _k, _v)

	return nil
}

// Delete deletes an instance of Person given its Address.
// Returns ErrPersonNotFound in case no record is found.
func (x PersonStateHandler) Delete(storage std.Storage, o Person) error {
	_k := keys.StringPrimaryKey(o.Address)
	_v := x.ns.Get(storage, _k)

	if _v == nil {
		return ErrPersonNotFound
	}

	x.ns.Remove(storage, _k)

	return nil
}
