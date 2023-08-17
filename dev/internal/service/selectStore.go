package service

import (
	"context"
	"errors"
	"fmt"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
)

func (ths *service) SelectStore(ctx context.Context, token *model.UserToken, req constant.SelectStoreRequest) error {
	isStoreAdmin, err := ths.repository.CheckUserIsStoreAdmin(ctx, token.UserID, req.ID)
	if err != nil {
		return fmt.Errorf("failed at repository.CheckUserIsStoreAdmin: %s", err.Error())
	}

	if isStoreAdmin {
		return errors.New("you are the admin of the store")
	}

	err = ths.repository.UpdateTokenStore(ctx, token.ID, req.ID)
	if err != nil {
		return fmt.Errorf("failed at repository.GetStoreByRef: %s", err.Error())
	}

	return nil
}
