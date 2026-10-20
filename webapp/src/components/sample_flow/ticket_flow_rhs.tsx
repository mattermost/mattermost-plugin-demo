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
import TicketList from './ticket_list';
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
                    <h3 className='DemoTicketFlowRhs__detailTitle'>{selected.key}</h3>
                </div>
                <div className='DemoTicketFlowRhs__scroll'>
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
                        emphasis='secondary'
                        size='small'
                        onClick={() => setModalOpen(true)}
                    >
                        <FormattedMessage
                            id='ticketFlow.rhs.openModal'
                            defaultMessage='Open in a modal'
                        />
                    </Button>
                    {pagePath && (
                        <Button
                            emphasis='secondary'
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
            {modalOpen && <TicketFlowModal onClose={() => setModalOpen(false)}/>}
        </div>
    );
}
