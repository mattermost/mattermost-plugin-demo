// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Channel} from 'mattermost-redux/types/channels';

function verifyExportBotMessage(channelDisplayName : string) : void {
    // # Find the last post in the current channel.
    cy.getLastPostId().then((lastPostId: string) => {
        // * Verify that the post contains a message about the export.
        cy.get(`#post_${lastPostId}`).
            should('contain.text', `Channel ~${channelDisplayName} exported:`);
    });
}
Cypress.Commands.add('verifyExportBotMessage', verifyExportBotMessage);

function verifyExportCommandIsAvailable() : void {
    // # Type the command in the textbox without sending the message.
    cy.findByTestId('post_textbox').clear().type('/export');

    // # Get the list of autocomplete suggestions.
    cy.get('#suggestionList').within(() => {
        // * Verify that there is a suggestion about the channel export command.
        cy.get('div.slash-command__desc').should('contain', 'Export the current channel');
    });
}
Cypress.Commands.add('verifyExportCommandIsAvailable', verifyExportCommandIsAvailable);

function verifyExportSystemMessage(channelDisplayName : string) : void {
    // # Find the last post in the current channel.
    cy.getLastPostId().then((lastPostId: string) => {
        // * Verify that the post contains a message about the export.
        cy.get(`#post_${lastPostId}`).
            should('contain.text',
                `Exporting ~${channelDisplayName}. @channelexport will send you a direct message when the export is ready.`);
    });
}
Cypress.Commands.add('verifyExportSystemMessage', verifyExportSystemMessage);

function verifyFileCanBeDownloaded(channelDisplayName : string) : void {
    const fileRegexp = new RegExp(Cypress.config('baseUrl') + '/api/v4/files/.*?download=1');

    // # Find the last post in the current channel.
    cy.getLastPostId().then((lastPostId: string) => {
        // * Verify that the post contains a message about the export.
        cy.get(`#post_${lastPostId}`).
            should('contain.text', `Channel ~${channelDisplayName} exported:`).
            within(() => {
                // # Find the list of files attached to the post.
                cy.findByTestId('fileAttachmentList').should('be.visible').within(() => {
                    // * Verify that the file has a download link with an href
                    // to a file in the system.
                    cy.get('a[download]').
                        should('have.attr', 'href').
                        should('match', fileRegexp);
                });
            });
    });
}
Cypress.Commands.add('verifyFileCanBeDownloaded', verifyFileCanBeDownloaded);

function verifyFileName(channelDisplayName : string, channelName: string) : void {
    // # Find the last post in the current channel.
    cy.getLastPostId().then((lastPostId: string) => {
        // * Verify that the post contains a message about the export.
        cy.get(`#post_${lastPostId}`).
            should('contain.text', `Channel ~${channelDisplayName} exported:`).
            within(() => {
                // # Find the list of files attached to the post.
                cy.findByTestId('fileAttachmentList').should('be.visible').within(() => {
                    // * Verify that the file to download has the name
                    // ${channelName}.csv
                    cy.get('a[download]').
                        should('have.attr', 'download', `${channelName}.csv`);
                });
            });
    });
}
Cypress.Commands.add('verifyFileName', verifyFileName);

function verifyNoPosts(channelName: string) : Cypress.Chainable<Channel> {
    // # Retrieve the channel via the API.
    return cy.apiGetChannelByName('ad-1', channelName).then((channel: Channel) => {
        // * Verify that there are no posts. There is always at least one post
        // in every channel (the system announcing the user joined) so we need
        // to check for 1 instead of 0.
        expect(channel.total_msg_count).at.most(1);
    });
}
Cypress.Commands.add('verifyNoPosts', verifyNoPosts);

function verifyAtLeastPosts(channelName: string, numPosts: number) : Cypress.Chainable<Channel> {
    // # Retrieve the channel via the API.
    return cy.apiGetChannelByName('ad-1', channelName).then((channel: Channel) => {
        // * Verify that the channel contains at least the number of messages specified.
        expect(channel.total_msg_count).at.least(numPosts);
    });
}
Cypress.Commands.add('verifyAtLeastPosts', verifyAtLeastPosts);
