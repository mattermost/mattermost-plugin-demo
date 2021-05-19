// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import * as TIMEOUTS from 'mattermost-webapp/e2e/cypress/fixtures/timeouts';

describe('Crash', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');

    before(() => {
        cy.apiAdminLogin();
        cy.visit('/');

        cy.apiRemovePluginById(pluginID, '');

        cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);
    });

    after(() => {
        cy.apiRemovePluginById(pluginID, '');
    });

    it('MM-T2404 crash and restart', () => {
        // # Post crash slash command
        cy.get('#post_textbox').clear().type('/crash {enter}');

        // * Verify ephemeral post confirming plugin crashes
        cy.get('#postListContent').should('contain.text', 'Crashing plugin');

        // # wait a few seconds for plugin to re-enable
        cy.wait(TIMEOUTS.TWO_SEC);

        // # Post plugin slash command to trigger plugin not working ephemeral response
        cy.get('#post_textbox').clear().type('/demo_plugin true {enter}');

        // * Verify ephemeral post confirming plugin crashes
        cy.findByText('Plugin for /demo_plugin is not working. Please contact your system administrator').should('be.visible');

        // # wait a few seconds for plugin to re-enable
        cy.wait(TIMEOUTS.HALF_MIN);

        // # @mention the demo plugin user
        cy.get('#post_textbox').clear().type('@demo_plugin hello {enter}');

        // * Confirm plugin is responsive again. Verify ephemeral message is posted
        cy.getLastPostId().then((postId) => {
            cy.get(`#post_${postId}`).
                findByText('Shh! You must not talk about the demo plugin user.').should('be.visible');
        });
    });
});
