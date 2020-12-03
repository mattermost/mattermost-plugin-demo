// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import * as TIMEOUTS from '../fixtures/timeouts';

function getLastPostId() : Cypress.Chainable<string> {
    // # Wait until permanent post.
    cy.get('#postListContent').should('be.visible');
    cy.waitUntil(() =>
        cy.findAllByTestId('postView').
            last().
            then((el) => !(el[0].id.includes(':')))
    );

    // # Get the last post and return its ID.
    return cy.findAllByTestId('postView').last().should('have.attr', 'id').and('not.include', ':').
        invoke('replace', 'post_', '');
}
Cypress.Commands.add('getLastPostId', getLastPostId);

// ***********************************************************
// Account Settings Modal
// ***********************************************************

// Go to Account Settings modal
Cypress.Commands.add('toAccountSettingsModal', () => {
    cy.get('#channel_view', {timeout: TIMEOUTS.ONE_MIN}).should('be.visible');
    cy.get('#sidebarHeaderDropdownButton').should('be.visible').click();
    cy.get('#accountSettings').should('be.visible').click();
    cy.get('#accountSettingsModal').should('be.visible');
});


// ***********************************************************
// Post
// ***********************************************************

Cypress.Commands.add('postMessage', (message) => {
    postMessageAndWait('#post_textbox', message);
});

Cypress.Commands.add('uiClickPostDropdownMenu', (postId, menuItem, location = 'CENTER') => {
    cy.clickPostDotMenu(postId, location);
    cy.findByTestId(`post-menu-${postId}`).should('be.visible');
    cy.findByText(menuItem).scrollIntoView().should('be.visible').click({force: true});
});

/**
 * Click dot menu by post ID or to most recent post (if post ID is not provided)
 * @param {String} postId - Post ID
 * @param {String} location - as 'CENTER', 'RHS_ROOT', 'RHS_COMMENT', 'SEARCH'
 */
Cypress.Commands.add('clickPostDotMenu', (postId, location = 'CENTER') => {
    clickPostHeaderItem(postId, location, 'button');
});

function postMessageAndWait(textboxSelector, message: string) {
    cy.get(textboxSelector, {timeout: TIMEOUTS.HALF_MIN}).should('be.visible').clear().type(`${message}{enter}`).wait(TIMEOUTS.HALF_SEC);
    cy.waitUntil(() => {
        return cy.get(textboxSelector).then((el) => {
            return el[0].textContent === '';
        });
    });
}

function clickPostHeaderItem(postId: string, location: string, item: string) {
    let idPrefix: string;
    switch (location) {
    case 'CENTER':
        idPrefix = 'post';
        break;
    case 'RHS_ROOT':
    case 'RHS_COMMENT':
        idPrefix = 'rhsPost';
        break;
    case 'SEARCH':
        idPrefix = 'searchResult';
        break;

    default:
        idPrefix = 'post';
    }

    if (postId) {
        cy.get(`#${idPrefix}_${postId}`).trigger('mouseover', {force: true});
        cy.wait(TIMEOUTS.HALF_SEC).get(`#${location}_${item}_${postId}`).click({force: true});
    } else {
        cy.getLastPostId().then((lastPostId) => {
            cy.get(`#${idPrefix}_${lastPostId}`).trigger('mouseover', {force: true});
            cy.wait(TIMEOUTS.HALF_SEC).get(`#${location}_${item}_${lastPostId}`).click({force: true});
        });
    }
}


