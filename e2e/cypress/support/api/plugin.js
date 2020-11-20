// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import * as TIMEOUTS from '../../fixtures/timeouts';

// *****************************************************************************
// Plugins
// https://api.mattermost.com/#tag/plugins
// *****************************************************************************

Cypress.Commands.add('apiUploadPlugin', (filename) => {
    return cy.apiUploadFile('plugin', filename, {url: '/api/v4/plugins', method: 'POST', successStatus: 201});
});

Cypress.Commands.add('apiEnablePluginById', (pluginId) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: `/api/v4/plugins/${encodeURIComponent(pluginId)}/enable`,
        method: 'POST',
        timeout: TIMEOUTS.ONE_MIN,
        failOnStatusCode: true,
    }).then((response) => {
        expect(response.status).to.equal(200);
        return cy.wrap(response);
    });
});

Cypress.Commands.add('apiRemovePluginById', (pluginId) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: `/api/v4/plugins/${encodeURIComponent(pluginId)}`,
        method: 'DELETE',
        failOnStatusCode: false,
    }).then((response) => {
        if (response.status !== 200 && response.status !== 404) {
            expect(response.status).to.equal(200);
        }

        return cy.wrap(response);
    });
});
