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

describe('Themes', () => {
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

    it('MM-T2403 theme', () => {
        // # change theme to mattermost dark
        navigateToThemeSettings();

        // * Verify icon color is white
        cy.get('.team-sidebar-bottom-plugin').find('.fa-plug').should('have.css', 'color', 'rgb(255, 255, 255)');
    });
});

function navigateToThemeSettings() {
    // # Change theme to desired theme (keeps settings modal open)
    cy.toAccountSettingsModal();
    cy.get('#displayButton').click();
    cy.get('#displaySettingsTitle').should('exist');

    // # Open edit theme
    cy.get('#themeTitle').should('be.visible');
    cy.get('#themeEdit').click();

    // # Click on the image
    cy.findByAltText('premade theme mattermostDark').click();
    cy.get('[data-testid="saveSetting"]').click();
    cy.get('#accountSettingsHeader').find('.close').click();
}

