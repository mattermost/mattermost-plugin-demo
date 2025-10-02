// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"

	"github.com/pkg/errors"
)

func (a *WhatsappApp) TryLock(key string, value []byte) (bool, error) {
	locked, appErr := a.api.KVCompareAndSet(key, nil, value)
	if appErr != nil {
		msg := fmt.Sprintf("TryLock: failed to save value in KV store via KVCompareAndSet, key: %s, value: %s, error: %s", key, value, appErr.Error())
		a.api.LogError(msg)
		return false, errors.New(msg)
	}

	return locked, nil
}

func (a *WhatsappApp) Unlock(key string, value []byte) (bool, error) {
	unlocked, appErr := a.api.KVCompareAndDelete(key, value)
	if appErr != nil {
		msg := fmt.Sprintf("Unlock: failed to delete KV store entry, key: %s, error: %s", key, appErr.Error())
		a.api.LogError(msg)
		return false, errors.New(msg)
	}

	return unlocked, nil
}
