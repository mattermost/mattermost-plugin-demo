// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export const httpStatusOk = 200;
export const httpStatusCreated = 201;
export const httpStatusNotFound = 404;

export const HTTPStatuses = {
    httpStatusOk: 200,
    httpStatusCreated: 201,
    httpStatusNotFound: 404
}

const MILLISECONDS_PER_SECOND = 1000;
const SECONDS_PER_MINUTE = 60;

const SECOND = MILLISECONDS_PER_SECOND;
const MINUTE = SECOND * SECONDS_PER_MINUTE;


export const TIMEOUTS = {
    ONE_HUNDRED_MILLIS: 100,
    HALF_SEC: SECOND / 2,
    ONE_SEC: SECOND,
    TWO_SEC: SECOND * 2,
    THREE_SEC: SECOND * 3,
    FOUR_SEC: SECOND * 4,
    FIVE_SEC: SECOND * 5,
    TEN_SEC: SECOND * 10,
    HALF_MIN: MINUTE / 2,
    ONE_MIN: MINUTE,
    TWO_MIN: MINUTE * 2,
    THREE_MIN: MINUTE * 3,
    FOUR_MIN: MINUTE * 4,
    FIVE_MIN: MINUTE * 5,
    TEN_MIN: MINUTE * 10,
    TWENTY_MIN: MINUTE * 20,
}

export const Constants = {
    HTTPStatuses,
    TIMEOUTS,
}

export default Constants

