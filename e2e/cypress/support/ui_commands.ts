// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
