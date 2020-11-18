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
         * Create a public channel via API.
         * @param {string} teamId - the ID of the team the channel will be created in.
         * @param {string} name - the name of the channel that will be created.
         * @param {string} displayName - the display name of the channel that will be created.
         * @return The created public channel.
        */
        apiCreatePublicChannel(teamId: string, name: string, displayName: string): Chainable<Channel>;

        /**
         * Create a private channel via API.
         * @param {string} teamId - the ID of the team the channel will be created in.
         * @param {string} name - the name of the channel that will be created.
         * @param {string} displayName - the display name of the channel that will be created.
         * @return The created private channel.
        */
        apiCreatePrivateChannel(teamId: string, name: string, displayName: string): Chainable<Channel>;

        /**
         * Create a group message (GM) with a set of users.
         * @param {string[]} userIds - an array of the IDs of the users to be added to the group.
         * @return The created group message channel
        */
        apiCreateGroupMessage(userIds: string[]): Chainable<Channel>;

        /**
         * Create a direct message (DM) with another user.
         * @param {string} selfId - ID of the creator of the DM.
         * @param {string} otherId - ID of the other user of the DM.
         * @return The created direct message channel
        */
        apiCreateDirectMessage(selfId: string, otherId: string): Chainable<Channel>;

        /**
         * Retrieve the information of a time given its name.
         * @param {string} name - the name of the team.
         * @return The retrieved team.
        */
        apiGetTeamByName(name: string): Chainable<Team>;

        /**
         * Retrieve the information of a user given its username.
         * @param {string} username - the username of the user.
         * @return The retrieved user.
        */
        apiGetUserByUsername(username: string): Chainable<UserProfile>;

        /**
         * Retrieve the information of a list of users given their usernames.
         * It's the same as apiGetUserByUsername, but with a list of users.
         * @param {string} username - the usernames of the users.
         * @return The list of retrieved users, in the same order than the ones
         * in the provided usernames list.
        */
        apiGetUsers(usernames: string[]): Chainable<UserProfile[]>;

        /**
         * Create a post via API.
         * @param {string} channelId - the ID of the channel where the message
         * will be posted.
         * @param {string} message - the message that will be posted.
         * @param {string[]} fieldIds - an optional list of IDs for attached
         * files; defaults an empty list.
         * @return Nothing.
        */
        apiCreatePost(channelId: string, message: string, fileIds?: string[]): Chainable<void>;

        /**
         * Retrieve the information of a channel given its team's name and its name.
         * @param {string} teamName - the name of the team the channel is in.
         * @param {string} channelName - the name of the channel.
         * @return The retrieved channel.
        */
        apiGetChannelByName(teamName: string, channelName: string) : Chainable<Channel>;

        /**
         * Make a channel, identified by its ID, read-only, meaning that
         * neither guests nor members can post any messages.
         * This can be done only by admins.
         * @param {string} channelId - the ID of the channel to make read-only.
         * @return The response of the API call.
        */
        apiMakeChannelReadOnly(channelId: string): Chainable<Response>;

        /**
         * Export a channel via the plugin's API endpoint.
         * @param {string} channelId - the ID of the channel to be exported.
         * @param {number} expectedStatus - the expected HTTP status code; it
         * defaults to 200.
         * @return The contents of the exported file.
        */
        apiExportChannel(channelId: string, expectedStatus?: number): Chainable<string>;

        /**
         * Retrieve the information of the logged in user.
         * @return The retrieved user.
        */
        apiGetMe(): Chainable<UserProfile>;

        /**
         * Update the admin configuration with the passed settings.
         * @param {PartialAdminConfig} newConfig - the partial configuration
         * that needs to be updated.
        */
        apiUpdateConfig(newConfig: PartialAdminConfig): Chainable<AdminConfig>;

        /**
         * Save a preference of a user directly via API.
         * This API assume that the user is logged in and has cookie to access.
         * @param {PreferenceType} preference - The user's preference to set.
         * @return The array of preferences of the user.
        */
        apiSaveUserPreference(preference: PreferenceType): Cypress.Chainable<PreferenceType[]>;

        /**
         * Save message display preference of a user directly via API.
         * This API assume that the user is logged in and has cookie to access.
         * @param {MessageDisplay} value - Either 'clean' or 'compact'.
         * @return The array of preferences of the user.
        */
        apiSaveMessageDisplayPreference(value: MessageDisplay) : Cypress.Chainable<PreferenceType[]>;

        /**
         * Save team mate name display preference of a user directly via API.
         * This API assume that the user is logged in and has cookie to access.
         * @param {TeammateNameFormat} value - Either 'username',
         * 'nickname_full_name' or 'full_name';
         * @return The array of preferences of the user.
        */
        apiSaveTeammateNameDisplayPreference(value: TeammateNameFormat) : Cypress.Chainable<PreferenceType[]>;

        /**
         * Verify that the server has a license.
         * @return Nothing.
        */
        apiRequireLicense() : Cypress.Chainable<void>;

        // *****************************************************************************
        //
        // UI Commands
        //
        // *****************************************************************************

        /**
         * Identify the last post in the current channel and find its ID.
         * @return the ID of the last post in the current channel.
        */
        getLastPostId(): Chainable<string>;

        /**
         * Post a message in the current channel.
         * @param {string} message - the string to post as a message in the current
         * channel.
         * @return Nothing.
        */
        postInCurrentChannel(message: string): Chainable<void>;

        /**
         * Send the /export slash command to the current channel.
         * @return Nothing.
        */
        exportSlashCommand(): Chainable<void>;

        /**
         * Archive the current channel via the channel header dropdown menu.
         * @return Nothing.
        */
        archiveCurrentChannel(): Chainable<void>;

        /**
         * Unarchive the current channel via the channel header dropdown menu.
         * @return Nothing.
        */
        unarchiveCurrentChannel(): Chainable<void>;

        /**
         * Leave the current channel via the channel header dropdown menu.
         * @return Nothing.
        */
        leaveCurrentChannel(): Chainable<void>;

        /**
         * Invite a user to the current channel via the channel header dropdown
         * menu.
         * @param {string} userName - the username of the user to be invited.
         * @return Nothing.
        */
        inviteUser(userName: string): Chainable<void>;

        /**
         * Kick a user from the current channel via the /kick slash command.
         * @param {string} userName - the username of the user to be kicked.
         * @return Nothing.
        */
        kickUser(userName: string): Chainable<void>;

        // *****************************************************************************
        //
        // Verification Commands
        //
        // *****************************************************************************

        /**
         * Verify that the bot sent a DM about the channel export.
         * @param {string} channelDisplayName - the expected name in the bot's
         * message: `Channel ~${channelDisplayName} exported`.
         * @return Nothing.
        */
        verifyExportBotMessage(channelDisplayName: string): Chainable<void>;

        /**
         * Verify that the export command is available, checking that it is
         * suggested by the autocompletion.
         * @return Nothing.
        */
        verifyExportCommandIsAvailable(): Chainable<void>;

        /**
         * Verify that the system sent an ephemeral message stating that the
         * channel export has started.
         * @param {string} channelDisplayName - the expected name in the bot's
         * message: `Exporting ~${channelDisplayName}. @channelexport will send
         * you a direct message when the export is ready.`
         * @return Nothing.
        */
        verifyExportSystemMessage(channelDisplayName: string): Chainable<void>;

        /**
         * Verify that the bot's DM contains a file to be downloaded.
         * @param {string} channelDisplayName - the expected name in the bot's
         * message: `Channel ~${channelDisplayName} exported`.
         * @return Nothing.
        */
        verifyFileCanBeDownloaded(channelDisplayName: string): Chainable<void>;

        /**
         * Verify that the bot's DM contains a file with the name `${channelName}.csv`
         * @param {string} channelDisplayName - the expected name in the bot's
         * message: `Channel ~${channelDisplayName} exported`.
         * @return Nothing.
        */
        verifyFileName(channelDisplayName: string, channelName: string): Chainable<void>;

        /**
         * Verify that the channel contains no posts (apart from the system one
         * informing that the user joined the channel).
         * @param channelName {string} - the name of the channel that should
         * have no posts.
         * @return The information of the channel.
        */
        verifyNoPosts(channelName: string): Chainable<Channel>;

        /**
         * Verify that the channel contains a minimum number of posts.
         * @param channelName {string} - the name of the channel.
         * @param numPosts {string} - the minimum number of posts that the
         * channel should have.
         * @return The information of the channel.
        */
        verifyAtLeastPosts(channelName: string, numPosts: number): Chainable<Channel>;

        // *****************************************************************************
        //
        // Navigation Commands
        //
        // *****************************************************************************

        /**
         * Create a new public channel and visit it.
         * @return The information of the created channel.
        */
        visitNewPublicChannel(): Chainable<Channel>;

        /**
         * Create a new private channel and visit it.
         * @return The information of the created channel.
        */
        visitNewPrivateChannel(): Chainable<Channel>;

        /**
         * Create a new group message (GM) with a list of users and visit
         * the new channel.
         * @param userNames {string[]} - the list of userNames from the users that
         * will be in the GM.
         * @return The information of the created channel.
        */
        visitNewGroupMessage(userNames: string[]): Chainable<Channel>;

        /**
         * Create a new direct message (DM) with another user and visit
         * the new channel.
         * @param creatorName {string} - the username of the creator of the DM.
         * @param otherName {string} - the username of the other user of the DM.
         * @return The information of the created channel.
        */
        visitNewDirectMessage(creatorName: string, otherName: string): Chainable<Channel>;

        /**
         * Visit the direct message (DM) channel between the logged in user and
         * the provided user.
         * @param userName {string} - the username of the other user in the DM.
         * @return Nothing.
        */
        visitDMWith(userName: string): Chainable<void>;

        /**
         * Visit the direct message (DM) channel with the channel export bot
         * the provided user.
         * @return Nothing.
        */
        visitDMWithBot(): Chainable<void>;
    }
}
