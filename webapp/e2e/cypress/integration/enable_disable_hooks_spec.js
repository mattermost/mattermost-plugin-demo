// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

describe('Enable and disable plugin hook events for Demo plugin', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');

    before(() => {        cy.visit('/');

        cy.apiRemovePluginById(pluginID);

        cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);
    });

    after(() => {
        cy.apiRemovePluginById(pluginID);
    });

    it('MM-T2411 Enable and disable plugin hook events for Demo plugin', () => {
        // # Type `/demo_plugin`
        cy.get('#post_textbox').clear().type('/demo_plugin');

        // * Verify autocomplete hint is correct
        cy.get('.slash-command__title').should('contain.text', 'demo_plugin');
        cy.get('.slash-command__desc').should('contain.text', 'Enables or disables the demo plugin hooks.');

        // # Type `/demo_plugin `
        cy.get('#post_textbox').type(' ');

        // * Check autocomplete options
        cy.get('#suggestionList > .slash-command:nth-child(2) .slash-command__title').should('contain.text', 'true');
        cy.get('#suggestionList > .slash-command:nth-child(2) .slash-command__desc').should('contain.text', 'Enable demo plugin hooks');
        cy.get('#suggestionList > .slash-command:nth-child(3) .slash-command__title').should('contain.text', 'false');
        cy.get('#suggestionList > .slash-command:nth-child(3) .slash-command__desc').should('contain.text', 'Disable demo plugin hooks');

        // # Choose "true"
        cy.get('#suggestionList > .slash-command:nth-child(2)').click();
        cy.get('#post_textbox').type('{enter}');

        // * Check for ephemeral post saying that hooks are enabled
        cy.getLastPostId().then((postId) => {
            cy.get(`#${postId}_message`).within(() => {
                cy.get('.post-message__text-container').should('contain.text', 'The demo plugin hooks are already enabled.');
            });
        });

        // * Check sidebar for hooks enabled
        cy.get('.sidebar--left__icons').should('contain.text', 'Demo Plugin: Enabled');

        // # Disable plugin hooks
        cy.get('#post_textbox').clear().type('/demo_plugin false');
        cy.get('#suggestionList > .slash-command:nth-child(1)').click();
        cy.get('#post_textbox').type('{enter}');

        // * Check for ephemeral post saying that hooks are disabled
        cy.getLastPostId().then((postId) => {
            cy.get(`#${postId}_message`).within(() => {
                cy.get('.post-message__text-container').should('contain.text', 'Disabled demo plugin hooks.');
            });
        });

        // * Check sidebar for hooks disabled
        cy.get('.sidebar--left__icons').should('contain.text', 'Demo Plugin: Disabled');

        // # Re-enable plugin hooks
        cy.get('#post_textbox').clear().type('/demo_plugin true');
        cy.get('#suggestionList > .slash-command:nth-child(1)').click();
        cy.get('#post_textbox').type('{enter}');

        // * Check for ephemeral post saying that hooks are enabled
        cy.getLastPostId().then((postId) => {
            cy.get(`#${postId}_message`).within(() => {
                cy.get('.post-message__text-container').should('contain.text', 'Enabled demo plugin hooks.');
            });
        });

        // * Check sidebar for hooks enabled
        cy.get('.sidebar--left__icons').should('contain.text', 'Demo Plugin: Enabled');
    });
});
