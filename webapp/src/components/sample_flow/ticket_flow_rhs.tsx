// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useId, useState} from 'react';
import {FormattedMessage, useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import ArrowLeftIcon from '@mattermost/compass-icons/components/arrow-left';
import {Button} from '@mattermost/compass-ui/components/button';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {IconButton} from '@mattermost/compass-ui/components/icon-button';

import messages from './messages';
import TicketFlowModal from './ticket_flow_modal';
import TicketForm from './ticket_form';
import TicketList, {Reported} from './ticket_list';
import {useTickets} from './tickets';

import './sample_flow.css';

type Props = {
    pagePath?: string;
};

export default function TicketFlowRhs({pagePath}: Props) {
    const {formatMessage} = useIntl();
    const tickets = useTickets();
    const formId = useId();
    const history = useHistory();
    const [selectedId, setSelectedId] = useState<string>();
    const [returnId, setReturnId] = useState<string>();
    const [modalOpen, setModalOpen] = useState(false);
    const [modalKey, setModalKey] = useState(0);
    const selected = tickets.find((ticket) => ticket.id === selectedId);

    const back = () => {
        setReturnId(selectedId);
        setSelectedId(undefined);
    };

    if (selected) {
        return (
            <div className='DemoTicketFlowRhs'>
                <div className='DemoTicketFlowRhs__detailHeader'>
                    <IconButton
                        aria-label={formatMessage(messages.backToTickets)}
                        icon={(
                            <Icon
                                glyph={<ArrowLeftIcon/>}
                                size='20'
                            />
                        )}
                        onClick={back}
                    />
                    <h3
                        className='DemoTicketFlowRhs__detailTitle'
                        title={selected.title}
                    >
                        {selected.title}
                    </h3>
                </div>
                <div className='DemoTicketFlowRhs__scroll'>
                    <p className='DemoTicketFlowRhs__detailMeta'>
                        {`${selected.key} • ${selected.project} • `}
                        <Reported ticket={selected}/>
                    </p>
                    <TicketForm
                        key={selected.id}
                        ticket={selected}
                        formId={formId}
                        onSaved={back}
                    />
                </div>
                <div className='DemoTicketFlowRhs__footer DemoTicketFlow__actions'>
                    <Button
                        type='submit'
                        form={formId}
                    >
                        <FormattedMessage {...messages.save}/>
                    </Button>
                    <Button
                        emphasis='tertiary'
                        onClick={back}
                    >
                        <FormattedMessage {...messages.cancel}/>
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className='DemoTicketFlowRhs'>
            <div className='DemoTicketFlowRhs__scroll'>
                <p className='DemoTicketFlowRhs__intro'>
                    <FormattedMessage
                        id='ticketFlow.rhs.intro'
                        defaultMessage='One list and one form, shown in three places. Edit a ticket here, or open the same flow in a modal or a full page.'
                    />
                </p>
                <div className='DemoTicketFlow__actions'>
                    <Button
                        emphasis='tertiary'
                        size='small'
                        onClick={() => {
                            setModalKey((key) => key + 1);
                            setModalOpen(true);
                        }}
                    >
                        <FormattedMessage
                            id='ticketFlow.rhs.openModal'
                            defaultMessage='Open in a modal'
                        />
                    </Button>
                    {pagePath && (
                        <Button
                            emphasis='tertiary'
                            size='small'
                            onClick={() => history.push(pagePath)}
                        >
                            <FormattedMessage
                                id='ticketFlow.rhs.openPage'
                                defaultMessage='Open full screen'
                            />
                        </Button>
                    )}
                </div>
                <section aria-labelledby={`${formId}-heading`}>
                    <h3
                        id={`${formId}-heading`}
                        className='DemoTicketFlow__heading'
                    >
                        <FormattedMessage {...messages.tickets}/>
                    </h3>
                    <div className='DemoTicketFlowRhs__list'>
                        <TicketList
                            tickets={tickets}
                            onSelect={setSelectedId}
                            focusId={returnId}
                        />
                    </div>
                </section>
            </div>
            <TicketFlowModal
                key={modalKey}
                open={modalOpen}
                onClose={() => setModalOpen(false)}
            />
        </div>
    );
}
