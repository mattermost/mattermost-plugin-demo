// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {HTTP, TIMEOUTS}  from '../constants';

// *****************************************************************************
// Plugins
// https://api.mattermost.com/#tag/plugins
// *****************************************************************************

Cypress.Commands.add('apiUploadPlugin', (filename) => {

    return cy.apiUploadFile('plugin', filename, {url: '/api/v4/plugins', method: 'POST', successStatus: HTTP.StatusCreated});
});

Cypress.Commands.add('apiEnablePluginById', (pluginId) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: `/api/v4/plugins/${encodeURIComponent(pluginId)}/enable`,
        method: 'POST',
        timeout: TIMEOUTS.ONE_MIN,
        failOnStatusCode: true,
    }).then((response) => {
        expect(response.status).to.equal(HTTP.StatusOk);
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
        if (response.status !== HTTP.StatusOk && response.status !== HTTP.StatusNotFound) {
            expect(response.status).to.equal(HTTP.StatusOk);
        }

        return cy.wrap(response);
    });
});
