// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {type FormEvent, useEffect, useId, useRef, useState} from 'react';
import {FormattedMessage, defineMessages, useIntl} from 'react-intl';

import {Combobox} from '@mattermost/compass-ui/components/combobox';
import {ErrorMessage} from '@mattermost/compass-ui/components/error-message';
import {Radio} from '@mattermost/compass-ui/components/radio';
import {SectionNotice} from '@mattermost/compass-ui/components/section-notice';
import {Select} from '@mattermost/compass-ui/components/select';
import {Switch} from '@mattermost/compass-ui/components/switch';
import {TextArea} from '@mattermost/compass-ui/components/text-area';
import {TextInput} from '@mattermost/compass-ui/components/text-input';

import {LABEL_OPTIONS, PRIORITY_OPTIONS, STATUS_OPTIONS, type Ticket, saveTicket} from './tickets';

const messages = defineMessages({
    errorTitle: {id: 'ticketFlow.form.errorTitle', defaultMessage: 'There are errors in the form below'},
    errorDescription: {id: 'ticketFlow.form.errorDescription', defaultMessage: 'Please resolve the errors before continuing.'},
    summary: {id: 'ticketFlow.form.summary', defaultMessage: 'Summary'},
    summaryRequired: {id: 'ticketFlow.form.summaryRequired', defaultMessage: 'Enter a summary'},
    description: {id: 'ticketFlow.form.description', defaultMessage: 'Description'},
    status: {id: 'ticketFlow.form.status', defaultMessage: 'Status'},
    labels: {id: 'ticketFlow.form.labels', defaultMessage: 'Labels'},
});

type Props = {
    ticket: Ticket;
    formId: string;
    onSaved: () => void;
};

export default function TicketForm({ticket, formId, onSaved}: Props) {
    const {formatMessage} = useIntl();
    const [draft, setDraft] = useState(ticket);
    const [submitted, setSubmitted] = useState(false);
    const summaryRef = useRef<HTMLInputElement>(null);
    const summaryErrorId = useId();

    const update = (patch: Partial<Ticket>) => setDraft((current) => ({...current, ...patch}));
    const summaryInvalid = submitted && !draft.title.trim();
    const statusOptions = STATUS_OPTIONS.map((option) => ({value: option.value, label: formatMessage(option.label)}));

    useEffect(() => {
        summaryRef.current?.focus();
    }, []);

    const submit = (e: FormEvent) => {
        e.preventDefault();
        setSubmitted(true);
        if (!draft.title.trim()) {
            summaryRef.current?.focus();
            return;
        }
        saveTicket({...draft, title: draft.title.trim()});
        onSaved();
    };

    return (
        <form
            id={formId}
            className='DemoTicketForm'
            onSubmit={submit}
            noValidate={true}
        >
            {summaryInvalid && (
                <SectionNotice
                    type='danger'
                    title={formatMessage(messages.errorTitle)}
                    description={formatMessage(messages.errorDescription)}
                />
            )}

            <section className='DemoTicketForm__group'>
                <h3 className='DemoTicketForm__groupTitle'>
                    <FormattedMessage
                        id='ticketFlow.form.details'
                        defaultMessage='Details'
                    />
                </h3>
                <div className='DemoTicketForm__field'>
                    <TextInput
                        ref={summaryRef}
                        label={formatMessage(messages.summary)}
                        value={draft.title}
                        onChange={(e) => update({title: e.target.value})}
                        invalid={summaryInvalid}
                        aria-invalid={summaryInvalid}
                        aria-describedby={summaryInvalid ? summaryErrorId : undefined}
                        maxLength={120}
                    />
                    {summaryInvalid && (
                        <ErrorMessage
                            id={summaryErrorId}
                            message={formatMessage(messages.summaryRequired)}
                        />
                    )}
                </div>
                <TextArea
                    label={formatMessage(messages.description)}
                    value={draft.description}
                    onChange={(e) => update({description: e.target.value})}
                    rows={3}
                />
            </section>

            <section className='DemoTicketForm__group'>
                <h3 className='DemoTicketForm__groupTitle'>
                    <FormattedMessage
                        id='ticketFlow.form.classification'
                        defaultMessage='Classification'
                    />
                </h3>
                <Select
                    label={formatMessage(messages.status)}
                    options={statusOptions}
                    value={draft.status}
                    onChange={(status) => update({status})}
                />
                <Combobox
                    label={formatMessage(messages.labels)}
                    multiple={true}
                    options={LABEL_OPTIONS}
                    value={draft.labels}
                    onChange={(labels) => update({labels: (labels as string[] | null) ?? []})}
                />
            </section>

            <section className='DemoTicketForm__group'>
                <h3 className='DemoTicketForm__groupTitle'>
                    <FormattedMessage
                        id='ticketFlow.form.tracking'
                        defaultMessage='Tracking'
                    />
                </h3>
                <fieldset className='DemoTicketForm__choices'>
                    <legend className='DemoTicketForm__choicesTitle'>
                        <FormattedMessage
                            id='ticketFlow.form.priority'
                            defaultMessage='Priority'
                        />
                    </legend>
                    <div className='DemoTicketForm__choiceGroup'>
                        {PRIORITY_OPTIONS.map((option) => (
                            <Radio
                                key={option.value}
                                name={`${formId}-priority`}
                                value={option.value}
                                checked={draft.priority === option.value}
                                onChange={() => update({priority: option.value})}
                            >
                                <FormattedMessage {...option.label}/>
                            </Radio>
                        ))}
                    </div>
                </fieldset>
                <Switch
                    checked={draft.notifyWatchers}
                    onChange={(e) => update({notifyWatchers: e.target.checked})}
                >
                    <FormattedMessage
                        id='ticketFlow.form.notifyWatchers'
                        defaultMessage='Notify watchers when this ticket changes'
                    />
                </Switch>
            </section>
        </form>
    );
}
