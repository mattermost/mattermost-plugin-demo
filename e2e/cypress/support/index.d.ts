// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

declare namespace Cypress {

    // We cannot use normal imports; otherwise, we would convert this to a
    // normal module, and we want it to be an ambient module in order to merge
    // our declaration of Cypress with the global one.
    // See https://stackoverflow.com/a/51114250/3248221 for more information
    type Channel = import('mattermost-redux/types/channels').Channel;
    type Team = import('mattermost-redux/types/teams').Team;
    type UserProfile = import('mattermost-redux/types/users').UserProfile;
    type AdminConfig = import('mattermost-redux/types/config').AdminConfig;
    type PreferenceType = import('mattermost-redux/types/preferences').PreferenceType;

    type PartialAdminConfig = import('./server_api_commands').PartialAdminConfig;
    type MessageDisplay = import('./server_api_commands').MessageDisplay;
    type TeammateNameFormat = import('./server_api_commands').TeammateNameFormat;

    interface Chainable<Subject> {

        // *****************************************************************************
        //
        // API commands (ttps://api.mattermost.com/)
        //
        // *****************************************************************************

        /**
         * User login directly via API.
         * @param {String} username - the username of the user.
         * @param {String} password - the password of the uesr.
         * @return The Response returned by the API call.
        */
        apiLogin(username?: string, password?: string | null): Chainable<Response>;

        /**
         * Identify the last post in the current channel and find its ID.
         * @return the ID of the last post in the current channel.
        */
        getLastPostId(): Chainable<string>;

        /**
         * Identify the last post in the current channel and find its ID.
         * @return the ID of the last post in the current channel.
        */
        toAccountSettingsModal(): Chainable<string>;
    }
}
