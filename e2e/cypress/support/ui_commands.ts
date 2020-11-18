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

function postInCurrentChannel(message: string): void {
    // # Type the message in the textbox and type Enter.
    cy.findByTestId('post_textbox').clear().type(message).type('{enter}');

    // * Verify that the textbox is empty after sending the message.
    cy.findByTestId('post_textbox').should('have.text', '');
}
Cypress.Commands.add('postInCurrentChannel', postInCurrentChannel);

function exportSlashCommand() : void {
    // # Post /export (with Esc to make sure that the autocompletion is gone)
    // in the channel.
    cy.postInCurrentChannel('/export {esc}');
}
Cypress.Commands.add('exportSlashCommand', exportSlashCommand);

function archiveCurrentChannel() : void {
    // # Click on the header dropdown arrow.
    cy.get('#channelHeaderDropdownIcon').click();

    // # Select the Archive Channel item in the menu.
    cy.get('#channelArchiveChannel').click();

    // # Confirm the archival of the channel clicking on the Archive button.
    cy.get('#deleteChannelModalDeleteButton').click();
}
Cypress.Commands.add('archiveCurrentChannel', archiveCurrentChannel);

function unarchiveCurrentChannel() : void {
    // # Click on the header dropdown arrow.
    cy.get('#channelHeaderDropdownIcon').click();

    // # Select the Unarchive Channel item in the menu.
    cy.get('#channelUnarchiveChannel').click();

    // # Confirm the unarchival of the channel clicking on the Unarchive button.
    cy.get('#unarchiveChannelModalDeleteButton').click();
}
Cypress.Commands.add('unarchiveCurrentChannel', unarchiveCurrentChannel);

function leaveCurrentChannel() : void {
    // # Click on the header dropdown arrow.
    cy.get('#channelHeaderDropdownIcon').click();

    // # Select the Leave Channel item in the menu.
    cy.get('#channelLeaveChannel').click();
}
Cypress.Commands.add('leaveCurrentChannel', leaveCurrentChannel);

function inviteUser(userName: string): void {
    // # Click on the header dropdown arrow.
    cy.get('#channelHeaderDropdownIcon').click();

    // # Select the Add Members item in the menu.
    cy.get('#channelAddMembers').click();

    // # Type the username in the search box.
    cy.get('#selectItems').type(userName);

    // # Click on the row containing the user.
    cy.get('#multiSelectList').within(() => {
        cy.findByText(`@${userName}`).click({force: true});
    });

    // # Click the Add button.
    cy.get('#saveItems').click();
    cy.get('#addUsersToChannelModal').should('not.exist');

    // * Make sure that the Add Members modal is gone.
    cy.get('#addUsersToChannelModal').should('not.exist');

    // * Verify that there is a system message informing that the user was added
    // to the channel.
    getLastPostId().then((lastPostId: string) => {
        cy.get(`#post_${lastPostId}`).
            should('contain.text', `@${userName} added to the channel by you.`);
    });
}
Cypress.Commands.add('inviteUser', inviteUser);

function kickUser(userName: string): void {
    // # Post /kick @{userName} (with Esc to make sure that the autocompletion
    // is gone) in the channel.
    cy.postInCurrentChannel(`/kick @${userName} {esc}`);

    // * Verify that there is a system message informing that the user was
    // kicked from the channel.
    getLastPostId().then((lastPostId: string) => {
        cy.get(`#post_${lastPostId}`).
            should('contain.text', `@${userName} was removed from the channel.`);
    });
}
Cypress.Commands.add('kickUser', kickUser);
