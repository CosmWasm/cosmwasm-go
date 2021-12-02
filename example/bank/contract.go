package bank

import (
	"github.com/cosmwasm/cosmwasm-go/std"
	"github.com/cosmwasm/cosmwasm-go/std/types"
)

type CreateDenom struct {
	Name string `json:"name"`
}

func executeCreateDenom(deps *std.Deps, msg CreateDenom) (types.Response, error) {
	err := DenomMetadataState.Create(deps.Storage, DenomMetadata{Name: msg.Name})
	if err != nil {
		return types.Response{}, err
	}

	return types.Response{
		Messages:   nil,
		Data:       nil,
		Attributes: nil,
		Events:     nil,
	}, nil
}

func executeUpdateDenom(deps *std.Deps, msg UpdateDenom) (types.Response, error) {
	// update denom
	err := DenomMetadataState.Update(deps.Storage, msg.Name, func(old *DenomMetadata) (*DenomMetadata, error) {
		if old == nil {
			return nil, ErrDenomMetadataNotFound
		}

		old.Decimal = msg.DecimalPlaces
		return old, nil
	})
	if err != nil {
		return types.Response{}, err
	}

	return types.Response{}, nil
}

func executeDeleteDenom(deps *std.Deps, msg DeleteDenom) (types.Response, error) {
	err := DenomMetadataState.Delete(deps.Storage, DenomMetadata{Name: msg.Name})
	if err != nil {
		return types.Response{}, err
	}

	return types.Response{}, nil
}

type UpdateDenom struct {
	DecimalPlaces int
	Name          string
}

type DeleteDenom struct {
	Name string
}

type ExecuteMsg struct {
	CreateDenom *CreateDenom `json:"create_denom"`
	UpdateDenom *UpdateDenom `json:"update_denom"`
	DeleteDenom *DeleteDenom `json:"delete_denom"`
}
