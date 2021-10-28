// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
/// <reference path="../support/index.d.ts" />

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

import * as TIMEOUTS from '../fixtures/timeouts';

import '@testing-library/cypress/add-commands';

/**
 * Note : This test requires the demo plugin tar file under fixtures folder.
 * Download :
 * make dist latest master and copy to ./e2e/cypress/fixtures/com.mattermost.demo-plugin-0.9.0.tar.gz
 */

describe('PostMenu', () => {
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

    it('MM-T3425 Post menu submenu items are selectable', () => {
      // # post a test message
      cy.postMessage('Hello')

      // # get last postID
      cy.getLastPostId().then((postId) => {

      // # click on dot menu
      cy.clickPostDotMenu(postId);

      cy.get('div.post-list__dynamic').should('be.visible').scrollTo('bottom', {duration: TIMEOUTS.ONE_SEC})

      // # click the Submenu Example menu
      cy.findByText('Submenu Example').trigger('mouseover')

      // * Verify submenu has 3 elements
      cy.get('.SubMenu').children().should('have.length', 3)

      // # click the first item submenu
      cy.findByText('First Item').trigger('mouseover').click()

      // * Verify root component shown with text
      cy.get('[data-testid="rootModalMessage"]').should('contain', 'You have triggered the root component of the demo plugin')
      });
    });

    it('MM-T2418-5 Post menu action', () => {
      // # post a test message
      cy.postMessage('Hello')

      // # get last postID
      cy.getLastPostId().then((postId) => {

        // # click on dot menu
        cy.clickPostDotMenu(postId)
      });

      // # scroll to the bottom so menu is visible
      cy.get('div.post-list__dynamic').should('be.visible').scrollTo('bottom', {duration: TIMEOUTS.ONE_SEC}).within(() => {

        // # click the Demo Plugin menu
        cy.findByText('Demo Plugin').trigger('mouseover').click()
      });


      // * Verify root component shown with text
      cy.get('[data-testid="rootModalMessage"]').should('contain', 'You have triggered the root component of the demo plugin')
    });
});
