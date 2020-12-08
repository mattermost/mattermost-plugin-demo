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

describe('Popover', () => {

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

    it('MM-T2418-2 Profile pop-over', () => {

      // # Visit a channel with another user
      cy.visit('/ad-1/channels/minus-6');

      // # click on user propfile to open pop-over
      cy.findByText('ruth.mason').should('be.visible').click();

      cy.get('#user-profile-popover').should('be.visible').within(() => {

        // * Verify user attributes text is visible
        cy.findByText('Demo Plugin: User Attributes').should('be.visible');

        // * Verify Action button is visible and click
        cy.findByText('Action').should('be.visible').
            trigger('mouseover').click();
      });

      // * Verify root component shown with text
      cy.findByTestId(`rootModalMessage`).should('contain', 'You have triggered the root component of the demo plugin')
    });
});



