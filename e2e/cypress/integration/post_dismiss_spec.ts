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

describe('Posts', () => {
    const pluginIdDemo = 'com.mattermost.demo-plugin';
    const demoFile = 'com.mattermost.demo-plugin-0.9.0.tar.gz';

    before(() => {
        cy.apiLogin('sysadmin');
        cy.visit('/');

        cy.apiRemovePluginById(pluginIdDemo, '');

        cy.apiUploadPlugin(demoFile);
        cy.apiEnablePluginById(pluginIdDemo);
    });

    after(() => {
        cy.apiRemovePluginById(pluginIdDemo, '');
    });

    it('MM-T2405 allow plugin to dismiss post', () => {
        // # at-mention the demo plugin user
        cy.get('#post_textbox').clear().type('@demo_plugin hello {enter}');

        // * Verify previously posted message is removed from center channel
        cy.findByText('@demo_plugin hello').should('not.be.visible');

        // * Verify ephemeral message is posted
        cy.findByText('Shh! You must not talk about the demo plugin user.').should('be.visible');
    });
});

