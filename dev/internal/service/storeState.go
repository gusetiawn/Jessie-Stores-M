package service

import (
	"context"
	"errors"
	"fmt"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

func (ths *service) StoreState(ctx context.Context, token *model.UserToken, req constant.StoreStateRequest) (*model.Store, error) {
	store, err := ths.repository.GetStoreByRef(ctx, req.Type, req.ReferenceID)
	if err != nil {
		return nil, fmt.Errorf("failed at repository.GetStoreByRef: %s", err.Error())
	}

	if store == nil {
		return nil, errors.New("store not found")
	}

	err = ths.repository.UpdateTokenStore(ctx, token.ID, store.ID)
	if err != nil {
		return nil, fmt.Errorf("failed at repository.UpdateTokenStore: %s", err.Error())
	}

	return store, nil
}
