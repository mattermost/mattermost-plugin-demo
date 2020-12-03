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

describe('FilePreview', () => {

    const pluginIdDemo = 'com.mattermost.demo-plugin'
    const demoFile = 'com.mattermost.demo-plugin-0.9.0.tar.gz';

    before(() => {
        cy.apiLogin('sysadmin')
        cy.visit('/');

        cy.apiRemovePluginById(pluginIdDemo, "");
    });

    after(() => {
        cy.apiRemovePluginById(pluginIdDemo, "");
    });

    it('MM-T3427 Demo Plugin FilePreviewOverride component', () => {
      // upload file.demo with message
      cy.get('#fileUploadInput').attachFile('file.demo');
      cy.get('#post_textbox').type('attach a file.demo {enter}')

      cy.getLastPostId().then((postId) => {

        // # click on file to display file preview
        cy.get(`#post_${postId}`).within(() => {
            cy.get('[data-testid="fileAttachmentList"]').contains('file.demo').click()
        })

        // * verify close button is not shown
        cy.get('.modal-image__background').should('contain', 'file.demo').should('not.contain', 'close')

        // # enable demo plugin to show FilePreviewOverride component
        cy.apiUploadPlugin(demoFile);
        cy.apiEnablePluginById(pluginIdDemo);

        // * verify close button is shown
        cy.get('.modal-image__background').should('contain', ['file.demo']).should('contain', 'Close')
      });
    });
});
