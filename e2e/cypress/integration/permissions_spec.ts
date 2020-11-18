// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
/// <reference path="../support/index.d.ts" />

import '@testing-library/cypress/add-commands';
import {Channel} from 'mattermost-redux/types/channels';
import {UserProfile} from 'mattermost-redux/types/users';

import {httpStatusNotFound} from '../support/constants';

describe('Test Area - Permissions', () => {
    const fileHeader =
    'Post Creation Time,User Id,User Email,User Type,User Name,Post Id,Parent Post Id,Post Message,Post Type';

    before(() => {
        // * Verify that the server is licensed, needed for all plugin features
        cy.apiRequireLicense();

        // # Login as sysadmin
        cy.apiLogin('sysadmin').then(() => {
            // # Enable ExperimentalViewArchivedChannels config as a sysadmin
            cy.apiUpdateConfig({
                TeamSettings: {
                    ExperimentalViewArchivedChannels: true,
                },
            });

            // # Set the teammate name format to username
            cy.apiSaveTeammateNameDisplayPreference('username');
        });

        // # Login as non-admin user.
        cy.apiLogin('user-1');

        // # Visit the default channel
        cy.visit('/');

        // # Set the expected user preferences:
        //   - Message Display to Standard
        //   - Teammate Name Format to username
        cy.apiSaveMessageDisplayPreference('clean');
        cy.apiSaveTeammateNameDisplayPreference('username');
    });

    beforeEach(() => {
        // # Login as non-admin user
        cy.apiLogin('user-1');

        // # Visit the default channel
        cy.visit('/');
    });

    it('ID 9 - User can export a public channel', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(channel.display_name);
        });
    });

    it('ID 10 - User can export a private channel', () => {
        // # Create a new private channel and visit it.
        cy.visitNewPrivateChannel().then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(channel.display_name);
        });
    });

    it('ID 11 - User can export a group message channel', () => {
        const userNames = ['user-1', 'aaron.medina', 'aaron.peterson'];

        // # Create a new group message and visit it.
        cy.visitNewGroupMessage(userNames).then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(channel.name);
        });
    });

    it('ID 12 - User can export a direct message channel', () => {
        // # Create a new direct message and visit it.
        cy.visitNewDirectMessage('user-1', 'anne.stone').then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(channel.name);
        });
    });

    it('ID 13 - User can export a direct message channel with self', () => {
        // # Get the user from the username to retrieve its ID.
        cy.apiGetUserByUsername('user-1').then((user: UserProfile) => {
            // # Visit the direct message channel with self.
            cy.visit('/ad-1/messages/@user-1');

            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(`${user.id}__${user.id}`);
        });
    });

    it('ID 14 - User can export a bot message channel', () => {
        // # Open the DM with the @channelexport bot.
        cy.visitDMWithBot();

        // # Run the /export slash command.
        cy.exportSlashCommand();

        // # Get the user from the username to retrieve its ID.
        cy.apiGetUserByUsername('user-1').then((user: UserProfile) => {
            // # Get the bot from its username to retrieve its ID.
            cy.apiGetUserByUsername('channelexport').then((bot: UserProfile) => {
                // * Verify that there is a file to be downloaded in the bot's
                // message.
                cy.verifyFileCanBeDownloaded(`${bot.id}__${user.id}`);
            });
        });
    });

    it('ID 15 - User can export archived channel', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Archive the current channel
            cy.archiveCurrentChannel();

            // * Verify that the channel can be exported via the plugin's API.
            cy.apiExportChannel(channel.id).then((fileContents: string) => {
                expect(fileContents).to.contain(fileHeader);
            });
        });
    });

    it('ID 16 - User can export an unarchived channel', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Archive the current channel
            cy.archiveCurrentChannel();

            // # Unarchive the channel logged in as an admin
            cy.apiLogin('sysadmin').then(() => {
                // # Visit the channel
                cy.visit(`/ad-1/channels/${channel.name}`);

                // * Verify that the channel has loaded
                cy.get('#channelHeaderTitle').contains(channel.display_name);

                // # Unarchive the channel
                cy.unarchiveCurrentChannel();
            });

            // # Visit the channel logged in as a regular user
            cy.apiLogin('user-1').then(() => {
                cy.visit(`/ad-1/channels/${channel.name}`);

                // # Run the /export slash command.
                cy.exportSlashCommand();

                // # Open the DM with the @channelexport bot.
                cy.visitDMWithBot();

                // * Verify that there is a file to be downloaded in the bot's
                // message.
                cy.verifyFileCanBeDownloaded(channel.display_name);
            });
        });
    });

    it('ID 15 (2nd one) - User cannot export a channel they are not added to', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel) => {
            // # Leave the current channel.
            cy.leaveCurrentChannel();

            // * Verify that the channel cannot be exported via the plugin's API.
            cy.apiExportChannel(channel.id, httpStatusNotFound);
        });
    });

    it('ID 17 - User cannot export a channel once they are ‘kicked’ from the channel', () => {
        let channel: Channel;

        // # Create a new channel as an admin and invite the user to it.
        cy.apiLogin('sysadmin').then(() => {
            cy.visitNewPublicChannel().then((newChannel) => {
                channel = newChannel;
                cy.inviteUser('user-1');
            });
        });

        // # Visit the channel as the regular user.
        cy.apiLogin('user-1').then(() => {
            cy.visit(`/ad-1/channels/${channel.name}`);

            // * Verify that the /export command is available.
            cy.verifyExportCommandIsAvailable();
        });

        // # Visit the channel as the admin and kick the user from it.
        cy.apiLogin('sysadmin').then(() => {
            cy.visit(`/ad-1/channels/${channel.name}`);
            cy.kickUser('user-1');
        });

        // # Visit the channel as the regular user.
        cy.apiLogin('user-1').then(() => {
            cy.visit('/ad-1/channels/town-square');

            // * Verify that the channel cannot be exported via the plugin's
            // API.
            cy.apiExportChannel(channel.id, httpStatusNotFound);
        });
    });

    it('ID 18 - User can export a read-only channel', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # As an admin, make the channel read-only.
            cy.apiLogin('sysadmin').then(() => {
                cy.apiMakeChannelReadOnly(channel.id);
            });

            // # Visit the channel as a regular user.
            cy.apiLogin('user-1').then(() => {
                cy.visit(`/ad-1/channels/${channel.name}`);

                // * Verify that the channel can be exported via the plugin's
                // API.
                cy.apiExportChannel(channel.id).then((fileContents: string) => {
                    expect(fileContents).to.contain(fileHeader);
                });
            });
        });
    });
});
