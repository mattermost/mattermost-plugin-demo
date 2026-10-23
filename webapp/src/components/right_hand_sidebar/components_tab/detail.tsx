import React, {useEffect, useRef, useState} from 'react';
import {defineMessages, useIntl} from 'react-intl';

import ArrowLeftIcon from '@mattermost/compass-icons/components/arrow-left';
import CheckIcon from '@mattermost/compass-icons/components/check';
import ContentCopyIcon from '@mattermost/compass-icons/components/content-copy';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {IconButton} from '@mattermost/compass-ui/components/icon-button';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';

import {importSnippet} from './helpers';
import type {CatalogEntry} from './types';

const COPIED_RESET_MS = 2000;

const messages = defineMessages({
    back: {id: 'components.detail.back', defaultMessage: 'Back to all components'},
    copy: {id: 'components.detail.copy', defaultMessage: 'Copy import'},
    copied: {id: 'components.detail.copied', defaultMessage: 'Copied'},
    copiedAnnouncement: {id: 'components.detail.copiedAnnouncement', defaultMessage: 'Import copied to clipboard'},
});

type Props = {
    entry: CatalogEntry;
    onBack: () => void;
};

export default function Detail({entry, onBack}: Props) {
    const {formatMessage} = useIntl();
    const DetailView = entry.detail;
    const backRef = useRef<HTMLButtonElement>(null);
    const [copied, setCopied] = useState(false);
    const snippet = importSnippet(entry);

    useEffect(() => {
        backRef.current?.focus({preventScroll: true});
    }, []);

    useEffect(() => {
        if (!copied) {
            return undefined;
        }
        const timeout = window.setTimeout(() => setCopied(false), COPIED_RESET_MS);
        return () => window.clearTimeout(timeout);
    }, [copied]);

    const copy = async () => {
        try {
            await navigator.clipboard.writeText(snippet);
            setCopied(true);
        } catch {
            setCopied(false);
        }
    };

    return (
        <div
            className='DemoCompassShowcase'
            onKeyDown={(event) => {
                if (event.key === 'Escape' && event.currentTarget.contains(event.target as Node)) {
                    event.stopPropagation();
                    onBack();
                }
            }}
        >
            <div className='DemoCompassShowcase__chrome DemoCompassShowcase__chrome--detail'>
                <div className='DemoCompassShowcase__detailHeader'>
                    <IconButton
                        ref={backRef}
                        aria-label={formatMessage(messages.back)}
                        icon={(
                            <Icon
                                glyph={<ArrowLeftIcon/>}
                                size='20'
                            />
                        )}
                        onClick={onBack}
                    />
                    <div className='DemoCompassShowcase__detailCopy'>
                        <h2 className='DemoCompassShowcase__detailTitle'>{entry.name}</h2>
                        <p className='DemoCompassShowcase__detailDescription'>{entry.description}</p>
                    </div>
                </div>
                <div className='DemoCompassShowcase__import'>
                    <code className='DemoCompassShowcase__importCode'>{snippet}</code>
                    <IconButton
                        aria-label={formatMessage(copied ? messages.copied : messages.copy)}
                        size='x-small'
                        icon={(
                            <Icon
                                glyph={copied ? <CheckIcon/> : <ContentCopyIcon/>}
                                size='16'
                            />
                        )}
                        onClick={copy}
                    />
                    <span
                        className='DemoCompassShowcase__srOnly'
                        aria-live='polite'
                    >
                        {copied ? formatMessage(messages.copiedAnnouncement) : ''}
                    </span>
                </div>
            </div>
            <Scrollbar className='DemoCompassShowcase__scroll'>
                <div className='DemoCompassShowcase__body DemoCompassShowcase__body--detail'>
                    <DetailView/>
                </div>
            </Scrollbar>
        </div>
    );
}
