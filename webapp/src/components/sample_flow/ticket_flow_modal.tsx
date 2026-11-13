// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useId, useState} from 'react';
import {FormattedMessage, useIntl} from 'react-intl';

import {Button} from '@mattermost/compass-ui/components/button';
import {Modal} from '@mattermost/compass-ui/components/modal';

import ModalLayer from '../modal_layer/modal_layer';

import messages from './messages';
import TicketForm from './ticket_form';
import TicketList, {Reported} from './ticket_list';
import {useTickets} from './tickets';

import './sample_flow.css';

type Props = {
    open: boolean;
    onClose: () => void;
};

export default function TicketFlowModal({open, onClose}: Props) {
    const {formatMessage} = useIntl();
    const tickets = useTickets();
    const formId = useId();
    const [selectedId, setSelectedId] = useState<string>();
    const [returnId, setReturnId] = useState<string>();
    const selected = tickets.find((ticket) => ticket.id === selectedId);

    const back = () => {
        setReturnId(selectedId);
        setSelectedId(undefined);
    };

    return (
        <ModalLayer
            open={open}
            onClose={onClose}
        >
            <Modal
                className='DemoTicketFlowModal'
                open={open}
                size='small'
                title={selected ? (
                    <span
                        className='DemoTicketFlowModal__title'
                        title={selected.title}
                    >
                        {selected.title}
                    </span>
                ) : formatMessage(messages.tickets)}
                subtitle={selected ? undefined : formatMessage(messages.pickTicket)}
                subtitlePlacement='beside'
                showBackButton={Boolean(selected)}
                onBack={back}
                backLabel={formatMessage(messages.backToTickets)}
                closeOnEscape={false}
                onClose={onClose}
                scrollable={true}
                footer={selected && (
                    <>
                        <Button
                            emphasis='tertiary'
                            onClick={back}
                        >
                            <FormattedMessage {...messages.cancel}/>
                        </Button>
                        <Button
                            type='submit'
                            form={formId}
                        >
                            <FormattedMessage {...messages.save}/>
                        </Button>
                    </>
                )}
            >
                {selected ? (
                    <>
                        <p className='DemoTicketFlowModal__meta'>
                            {`${selected.key} • ${selected.project} • `}
                            <Reported ticket={selected}/>
                        </p>
                        <TicketForm
                            key={selected.id}
                            ticket={selected}
                            formId={formId}
                            onSaved={back}
                        />
                    </>
                ) : (
                    <div className='DemoTicketFlowModal__list'>
                        <TicketList
                            tickets={tickets}
                            onSelect={setSelectedId}
                            focusId={returnId}
                        />
                    </div>
                )}
            </Modal>
        </ModalLayer>
    );
}
