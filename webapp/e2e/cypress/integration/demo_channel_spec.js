// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import * as TIMEOUTS from 'mattermost-webapp/e2e/cypress/fixtures/timeouts';

describe('Plugin functions - Demo channel', () => {
    const pluginID = Cypress.config('pluginID');
    const pluginFile = Cypress.config('pluginFile');
    let testTeam;

    before(() => {
        cy.apiInitSetup().then(({team}) => {
            testTeam = team;
        });
        // # Login as sysadmin
        cy.apiAdminLogin();

        cy.apiRemovePluginById(pluginID);

        cy.apiUploadPlugin(pluginFile);
        cy.apiEnablePluginById(pluginID);
    });

    after(() => {
        cy.apiRemovePluginById(pluginID);
    });

    it('MM-T2406 Plugin functions - Demo channel', () => {
        // # Visit bot account page
        cy.visit(`/${testTeam.name}/integrations/bots`);
        cy.get('.backstage-list__item .item-details__name', {timeout: TIMEOUTS.ONE_MIN}).should('be.visible');

        // * Verify demo plugin has created a bot account
        cy.get('.backstage-list__item .item-details__name').should('contain.text', 'Demo Plugin Bot (@demoplugin)');
        cy.get('.backstage-list__item .bot-details__description').should('contain.text', 'A bot account created by the demo plugin.');
        cy.get('.backstage-list__item .small').should('contain.text', 'Managed by plugin');

        // # Visit demo plugin channel
        cy.visit(`/${testTeam.name}/channels/demo_plugin`);

        // * Verify the plugin's bot has created a post
        cy.get('.post--bot.other--root', {timeout: TIMEOUTS.ONE_MIN}).should('contain.text', 'demoplugin');
        cy.get('.post--bot.other--root').should('contain.text', 'BOT');
    });
});
