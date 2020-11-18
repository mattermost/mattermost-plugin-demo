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

describe('Integrations', () => {
    const pluginIdDemo = 'com.mattermost.demo-plugin';

    before(() => {
        // # Initialize setup and visit town-square

        // * Verify that the server is licensed, needed for all plugin features
        cy.apiRequireLicense();

        // # Login as sysadmin
        // cy.apiLogin('sysadmin')
        cy.apiLogin('user-1')

        // # Visit the default channel.
        cy.visit('/');

        // # Set the expected user preferences:
        //   - Message Display to Standard
        //   - Teammate Name Format to username
        // cy.apiSaveMessageDisplayPreference('clean');
        // cy.apiSaveTeammateNameDisplayPreference('username');

    });

    after(() => {
        // cy.apiRemovePluginById(pluginIdDemo);
    });

    it('PASSES - MM-T2405 allow plugin to dismiss post', () => {
        // # at-mention the demo plugin user
        cy.get('#post_textbox').clear().type('@demo_plugin hello {enter}');

        // * Verify previously posted message is removed from center channel
        cy.findByText('@demo_plugin hello').should('not.be.visible');

        // * Verify ephemeral message is posted
        cy.findByText('Shh! You must not talk about the demo plugin user.').should('be.visible');
    });

    it('PASSES - MM-T2404 crash and restart', () => {
        // # Post crash slash command
        cy.get('#post_textbox').clear().type('/crash {enter}');

        // * Verify ephemeral post confirming plugin crashes
        cy.get('#postListContent').should('contain.text', 'Crashing plugin');

        // # Post crash slash command
        cy.get('#post_textbox').clear().type('/demo_plugin true {enter}');

        // * Verify ephemeral post confirming plugin crashes
        cy.findByText('Plugin for /demo_plugin is not working. Please contact your system administrator').should('be.visible');

        // # wait a few seconds for plugin to re-enable
        cy.wait(TIMEOUTS.HALF_MIN);

        // # @mention the demo plugin user
        cy.get('#post_textbox').clear().type('@demo_plugin hellow {enter}');

        // * Confirm plugin is responsive again. Verify ephemeral message is posted
        cy.findByText('Shh! You must not talk about the demo plugin user.').should('be.visible');
    });

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

    it('MM-T3422 Demo plugin can draw a tooltip', () => {
        // # Post a slash command that omits the optional argument
        cy.get('#post_textbox').clear().type('www.test.com {enter}');

        cy.getLastPostId().then((postId) => {
          cy.get(`#post_${postId}`).
            findByText('www.test.com').
            trigger('mouseover');
        });

        cy.findByText('BROKEN - cannot find text below').should('be.visible')
        cy.findByText('This is a custom tooltip').should('be.visible')
        // cy.get('div.span').should('be.visible').and('contain', 'This is a ');


        // // * Verify item is sent and is (only visible to you)
        // cy.getLastPostId().then((postId) => {
        //     cy.get(`#post_${postId}`).contains('Executed command: /autocomplete_test optional-arg');
        // });
    });

    it.only('MM-T3426 Hashtags still work with demo plugin enabled', () => {
        // # Post a
        cy.get('#post_textbox').clear().type('#pickles {enter}');

        let postID1
        // * Verify suggestion list is visible with 11 children
        cy.getLastPostId().then((postId: string) => {
          postID1 = postId
          cy.get(`#postMessageText_${postId}`).
            find('.mention-link').
            click();
        });

        // cy.get('#search-items-container').should('be.visible').within(() => {
            cy.get(`#searchResult_${postID1}`).
                should('contain.text', '#pickles')
        // });

    });
});

