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

describe('DemoChannel', () => {

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

    it('MM-T2412 Demo plugin rejects posts in the demo channel', () => {
      // # visit demo_channel
      cy.visit(`/ad-1/channels/demo_plugin`);

      // # post a test message
      cy.postMessage('Hello')

      // * verify posting in channel is not permitted
      cy.getLastPostId().then((postId) => {
        cy.get(`#post_${postId}`).should('contain', 'Posting is not allowed in this channel')
      })

      // # visit another channel
      cy.visit(`/ad-1/channels/town-square`);

      // # post a test message
      cy.postMessage('hello, @demo_plugin')

      // * verify return message
      cy.getLastPostId().then((postId) => {
        cy.get(`#post_${postId}`).should('contain', 'You must not talk about the demo plugin user')
      })

    });
});
