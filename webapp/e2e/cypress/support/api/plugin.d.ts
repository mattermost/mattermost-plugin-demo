// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/// <reference types="cypress" />

// ***************************************************************
// Each command should be properly documented using JSDoc.
// See https://jsdoc.app/index.html for reference.
// Basic requirements for documentation are the following:
// - Meaningful description
// - Specific link to https://api.mattermost.com
// - Each parameter with `@params`
// - Return value with `@returns`
// - Example usage with `@example`
// Custom command should follow naming convention of having `api` prefix, e.g. `apiLogin`.
// ***************************************************************

declare namespace Cypress {
    interface Chainable<Subject = any> {

        /**
         * Upload plugin.
         * See https://api.mattermost.com/#tag/plugins/paths/~1plugins/post
         * @param {string} filename - name of the plugin to upload
         * @returns {Response} response: Cypress-chainable response
         *
         * @example
         *   cy.apiUploadPlugin('filename');
         */
        apiUploadPlugin(filename: string): Chainable<Response>;

        /**
         * Enable plugin.
         * See https://api.mattermost.com/#tag/plugins/paths/~1plugins~1{plugin_id}~1enable/post
         * @param {string} pluginId - Id of the plugin to enable
         * @returns {string} `out.status`
         *
         * @example
         *   cy.apiEnablePluginById('pluginId');
         */
        apiEnablePluginById(pluginId: string): Chainable<Record<string, any>>;

        /**
         * Remove plugin.
         * See https://api.mattermost.com/#tag/plugins/paths/~1plugins~1{plugin_id}/delete
         * @param {string} pluginId - Id of the plugin to uninstall
         * @returns {string} `out.status`
         *
         * @example
         *   cy.apiRemovePluginById('url');
         */
        apiRemovePluginById(pluginId: string, force: string): Chainable<Record<string, any>>;
    }
}
