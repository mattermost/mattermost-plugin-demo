// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

func (a *WhatsappApp) ResetData() error {
	return a.store.ResetData()
}
