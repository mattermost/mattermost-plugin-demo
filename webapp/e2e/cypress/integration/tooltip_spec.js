// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

describe('Tooltips', () => {
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

    it('MM-T3422 Demo plugin can draw a tooltip', () => {
        // # Post a slash command that omits the optional argument
        cy.get('#post_textbox').clear().type('www.test.com {enter}');

        cy.getLastPostId().then((postId) => {
            cy.get(`#post_${postId}`).
                findByText('www.test.com').
                trigger('mouseover');
        });

        // * Verify tooltip show
        cy.findByTestId('tooltipMessage').should('have.text', 'This is a custom tooltip from the Demo Plugin');
    });
});

