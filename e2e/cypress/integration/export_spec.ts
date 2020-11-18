// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
/// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import '@testing-library/cypress/add-commands';
import {Channel} from 'mattermost-redux/types/channels';

describe('Test Area - Export', () => {
    before(() => {
        // * Verify that the server is licensed, needed for all plugin features
        cy.apiRequireLicense();

        // # Login as non-admin user.
        cy.apiLogin('user-1');

        // # Visit the default channel.
        cy.visit('/');

        // # Set the expected user preferences:
        //   - Message Display to Standard
        //   - Teammate Name Format to username
        cy.apiSaveMessageDisplayPreference('clean');
        cy.apiSaveTeammateNameDisplayPreference('username');
    });

    beforeEach(() => {
        // # Login as non-admin user.
        cy.apiLogin('user-1');

        // # Visit the default channel.
        cy.visit('/');
    });

    it('ID 19 - A system message notifies of successful export command execution on the channel where export is initiated', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // * Verify that there is a system message informing on the
            // export.
            cy.verifyExportSystemMessage(channel.display_name);
        });
    });

    it('ID 20 - A bot message notifies of a successful export', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that the bot sent a message informing about the export
            // of the specified channel.
            cy.verifyExportBotMessage(channel.display_name);
        });
    });

    it('ID 21 - The exported file can be downloaded locally', () => {
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

    it('ID 23 - Exported CSV filename has [channel-name].csv format', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that the file sent by the bot has the name
            // channelName.csv.
            cy.verifyFileName(channel.display_name, channel.name);
        });
    });

    it('ID 29 - A channel with no messages can be exported successfully', () => {
        // # Create a new public channel and visit it.
        cy.visitNewPublicChannel().then((channel: Channel) => {
            // * Make sure that there are no posts in the channel.
            cy.verifyNoPosts(channel.name);

            // # Run the /export slash command.
            cy.exportSlashCommand();

            // # Open the DM with the @channelexport bot.
            cy.visitDMWithBot();

            // * Verify that there is a file to be downloaded in the bot's
            // message.
            cy.verifyFileCanBeDownloaded(channel.display_name);
        });
    });

    it('ID 30 - A channel with more than 100 messages can be exported successfully', () => {
        const minPosts = 101;

        // * Make sure that the minima-3 channel contains at least 100 messages.
        cy.verifyAtLeastPosts('minima-3', minPosts).then((channel: Channel) => {
            // # Visit the minima-3 channel
            cy.visit('/ad-1/channels/minima-3');

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
