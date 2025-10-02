// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import "fmt"

const (
	UserLockKeyPrefix = "user_lock_"
)

func KeyUserSurveySentStatus(userID, surveyID string) string {
	return fmt.Sprintf("user_survey_status_%s_%s", userID, surveyID)
}

func KeyUserTeamMembershipFilterCache(userID, surveyID string) string {
	return fmt.Sprintf("user_team_filter_cache_%s_%s", userID, surveyID)
}

func KeyUserSendSurveyLock(userID string) string {
	return UserLockKeyPrefix + userID
}
