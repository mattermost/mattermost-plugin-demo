// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import users from '../fixtures/users';
import {httpStatusOk} from '../support/constants';

function apiLogin(username = 'user-1', password : string | null = null) : Cypress.Chainable<Cypress.Response> {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: '/api/v4/users/login',
        method: 'POST',
        body: {
            login_id: username,
            password: password || users[username].password,
        },
    }).then((response: Cypress.Response) => {
        expect(response.status).to.equal(httpStatusOk);
        return cy.wrap(response);
    });
}
Cypress.Commands.add('apiLogin', apiLogin);

Cypress.Commands.add('apiLogout', () => {
    cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: '/api/v4/users/logout',
        method: 'POST',
        log: false,
    });

    // * Verify logged out
    cy.visit('/login?extra=expired').url().should('include', '/login');

    // # Ensure we clear out these specific cookies
    ['MMAUTHTOKEN', 'MMUSERID', 'MMCSRF'].forEach((cookie) => {
        cy.clearCookie(cookie);
    });

    // # Clear remainder of cookies
    cy.clearCookies();

    // * Verify cookies are empty
    cy.getCookies({log: false}).should('be.empty');
});


