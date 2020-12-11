// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
interface User {
    username: string;
    password: string;
    email: string;
    emailname?: string;
    firstName: string;
    lastName: string;
}

const users: { [name: string]: User } = {
    sysadmin: {
        username: 'sysadmin',
        password: 'Sys@dmin-sample1',
        email: 'sysadmin@sample.mattermost.com',
        firstName: 'Kenneth',
        lastName: 'Moreno',
    },
    'user-1': {
        username: 'user-1',
        password: 'SampleUs@r-1',
        email: 'user-1@sample.mattermost.com',
        firstName: 'Victor',
        lastName: 'Welch',
    },
    'user-2': {
        username: 'samuel.tucker',
        password: 'SampleUs@r-2',
        email: 'user-2@sample.mattermost.com',
        emailname: 'user-2',
        firstName: 'Samuel',
        lastName: 'Tucker',
    },
};

export default users;
