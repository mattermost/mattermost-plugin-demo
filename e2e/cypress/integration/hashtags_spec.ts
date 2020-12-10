// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
/// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import '@testing-library/cypress/add-commands';

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

describe('Hashtags', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');

    before(() => {
        cy.apiLogin('sysadmin');
        cy.visit('/');

        cy.apiRemovePluginById(pluginID, '');

        cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);
    });

    after(() => {
        cy.apiRemovePluginById(pluginID, '');
    });

    it('MM-T3426 Hashtags still work with demo plugin enabled', () => {
        // # Post a hashtag
        cy.get('#post_textbox').clear().type('#pickles {enter}');

        cy.getLastPostId().then((postId: string) => {
            // # click hashtag in from the last post
            cy.get(`#postMessageText_${postId}`).
                find('.mention-link').
                click();

            // * verify post exists in the RHS
            cy.get(`#searchResult_${postId}`).
                should('contain.text', '#pickles');
        });
    });
});

