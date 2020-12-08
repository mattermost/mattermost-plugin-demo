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
         * Logout a user's active session from server via API.
         * See https://api.mattermost.com/#tag/users/paths/~1users~1logout/post
         * Clears all cookies espececially `MMAUTHTOKEN`, `MMUSERID` and `MMCSRF`.
         *
         * @example
         *   cy.apiLogout();
         */
        apiLogout(): Chainable<Response>;

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

        /**
         * Identify the last post in the current channel and find its ID.
         * @return the ID of the last post in the current channel.
        */
        postMessage(message: string | null): Chainable<string>;

        /**
         * Click dot menu by post ID or to most recent post (if post ID is not provided)
         * @param {String} postId - Post ID
         * @param {String} location - as 'CENTER', 'RHS_ROOT', 'RHS_COMMENT', 'SEARCH'
         */
        clickPostDotMenu(postId: string): Chainable<string>;

        /**
         * Click comment icon by post ID or to most recent post (if post ID is not provided)
         * This open up the RHS
         * @param {String} postId - Post ID
         * @param {String} location - as 'CENTER', 'SEARCH'
         */
        clickPostCommentIcon(postId: string): Chainable<string>;
    }
}
