// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
/// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import '@testing-library/cypress/add-commands';
import {Channel} from 'mattermost-redux/types/channels';

import * as TIMEOUTS from '../fixtures/timeouts';

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

describe('Integrations', () => {

    const pluginIdDemo = 'com.mattermost.demo-plugin'
    const demoFile = 'com.mattermost.demo-plugin-0.9.0.tar.gz';

    before(() => {
        // * Verify that the server is licensed, needed for all plugin features
        cy.apiRequireLicense();


        // # Login as sysadmin
        cy.apiLogin('sysadmin')
        cy.visit('/');

        cy.apiRemovePluginById(pluginIdDemo, "");

        cy.apiUploadPlugin(demoFile);
        cy.apiEnablePluginById(pluginIdDemo);

        // cy.apiInstallPluginFromUrl(demoURL, true);
    });

    beforeEach(() => {
        // cy.apiLogin('user-1')

        // # Visit the default channel.
        cy.apiLogin('sysadmin')
        cy.visit('/');

        // # Set the expected user preferences:
        //   - Message Display to Standard
        //   - Teammate Name Format to username
        // cy.apiSaveMessageDisplayPreference('clean');
        // cy.apiSaveTeammateNameDisplayPreference('username');
    });

    // after(() => {
        // cy.apiRemovePluginById(pluginIdDemo, "");
    // });

    it.skip('MM-T2403 theme', () => {
        // # Post a slash command with trailing space
        // # Post a slash command with trailing space
        cy.get('#post_textbox').clear().type('/autocomplete_test dynamic-arg ');

        // * Verify suggestion list is visible with at three children (issue, instance, info)
        cy.get('#suggestionList').should('be.visible').children().
            should('contain.text', 'suggestion 1 (hint)').
            should('contain.text', 'suggestion 2 (hint)');

        // # down arrow to highlight the second suggestion
        // # enter to push to send command to post textbox
        cy.get('#post_textbox').type('{downarrow}{downarrow}{enter}');

        // # Send the command
        cy.get('#post_textbox').type('{enter}');

        // * Verify correct message is sent,  (only visible to you)
        cy.getLastPostId().then((postId) => {
            cy.get(`#post_${postId}`).
                should('contain.text', '(Only visible to you)').
                should('contain.text', 'Executed command: /autocomplete_test dynamic-arg suggestion 2');
        });
    });
});

