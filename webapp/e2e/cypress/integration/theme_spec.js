// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

describe('Themes', () => {
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

    it('MM-T2403 theme', () => {
        cy.toAccountSettingsModal();

        // # Navigate to Display settings
        cy.get('#displayButton').click();
        cy.get('#displaySettingsTitle').should('exist');

        // # Open edit theme settings
        cy.get('#themeTitle').should('be.visible');
        cy.get('#themeEdit').click();

        // # Set theme to Mattermost Dark theme
        cy.findByAltText('premade theme mattermostDark').click();
        cy.get('[data-testid="saveSetting"]').click();
        cy.get('#accountSettingsHeader').find('.close').click();

        // * Verify icon color is white
        cy.get('.team-sidebar-bottom-plugin').find('.fa-plug').should('have.css', 'color', 'rgb(255, 255, 255)');
    });
});

