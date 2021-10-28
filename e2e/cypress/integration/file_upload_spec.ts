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

describe('File Upload', () => {

    const pluginIdDemo = 'com.mattermost.demo-plugin'
    const demoFile = 'com.mattermost.demo-plugin-0.9.0.tar.gz';

    before(() => {
        cy.apiLogin('sysadmin')
        cy.visit('/');
        cy.apiRemovePluginById(pluginIdDemo, "");

        cy.apiUploadPlugin(demoFile);
        cy.apiEnablePluginById(pluginIdDemo);
    });

    after(() => {
        cy.apiRemovePluginById(pluginIdDemo, "");
    });

    it('MM-T2418-7 File Upload Action', () => {

      // # click the file attachment in the main input box
      cy.get('#fileUploadButton').click()

      // * Verify demo plugin upload option shows
      cy.get('.Menu__content').should('contain', 'Upload using Demo Plugin').trigger('mouseover').click()

      // * Verify root component shown with text
      cy.findByTestId(`rootModalMessage`).should('contain', 'You have triggered the root component of the demo plugin').click()

      // # post a test message
      cy.postMessage('Hello')

      // # open RHS using previous post
      cy.getLastPostId().then((postId) => {
        // # click on dot menu
        cy.clickPostCommentIcon(postId);
      })

      cy.get(`#rhsContainer`).should('be.visible').within(() => {
        // # click the file attachment in the main input box
        cy.get('#fileUploadButton').click()

        // * Verify demo plugin upload option shows
        cy.get('.Menu__content').should('contain', 'Upload using Demo Plugin').trigger('mouseover').click()

      })

      // * Verify root component shown with text
      cy.findByTestId(`rootModalMessage`).should('contain', 'You have triggered the root component of the demo plugin')
    });
});



