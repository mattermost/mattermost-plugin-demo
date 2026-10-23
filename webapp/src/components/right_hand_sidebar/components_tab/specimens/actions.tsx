// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import CheckIcon from '@mattermost/compass-icons/components/check';
import PlusIcon from '@mattermost/compass-icons/components/plus';
import SendIcon from '@mattermost/compass-icons/components/send';
import {ActionButton} from '@mattermost/compass-ui/components/action-button';
import {Button} from '@mattermost/compass-ui/components/button';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {IconButton} from '@mattermost/compass-ui/components/icon-button';

import {Row, Variant} from './shared';

export function ButtonPreview() {
    return (
        <Row>
            <Button>{'Save'}</Button>
            <Button emphasis='secondary'>{'Cancel'}</Button>
        </Row>
    );
}

export function ButtonDetail() {
    return (
        <>
            <Variant label='Emphasis'>
                <Row>
                    <Button emphasis='primary'>{'Primary'}</Button>
                    <Button emphasis='secondary'>{'Secondary'}</Button>
                    <Button emphasis='tertiary'>{'Tertiary'}</Button>
                    <Button emphasis='quaternary'>{'Quaternary'}</Button>
                </Row>
            </Variant>
            <Variant label='Size'>
                <Row>
                    <Button size='x-small'>{'X-Small'}</Button>
                    <Button size='small'>{'Small'}</Button>
                    <Button size='medium'>{'Medium'}</Button>
                    <Button size='large'>{'Large'}</Button>
                </Row>
            </Variant>
            <Variant label='Destructive'>
                <Row>
                    <Button destructive={true}>{'Delete'}</Button>
                    <Button
                        disabled={true}
                        destructive={true}
                    >
                        {'Disabled'}
                    </Button>
                </Row>
            </Variant>
            <Variant label='With icon'>
                <Button
                    leadingIcon={(
                        <Icon
                            glyph={<SendIcon/>}
                            size='16'
                        />
                    )}
                >
                    {'Send'}
                </Button>
            </Variant>
        </>
    );
}

export function IconButtonPreview() {
    return (
        <IconButton
            aria-label='Add'
            icon={(
                <Icon
                    glyph={<PlusIcon/>}
                    size='20'
                />
            )}
        />
    );
}

export function IconButtonDetail() {
    const icon = (
        <Icon
            glyph={<PlusIcon/>}
            size='20'
        />
    );
    return (
        <>
            <Variant label='Default / toggled / destructive'>
                <Row>
                    <IconButton
                        aria-label='Add'
                        icon={icon}
                    />
                    <IconButton
                        aria-label='Selected'
                        icon={icon}
                        toggled={true}
                    />
                    <IconButton
                        aria-label='Delete'
                        destructive={true}
                        icon={icon}
                    />
                </Row>
            </Variant>
            <Variant label='Size'>
                <Row>
                    <IconButton
                        aria-label='X-Small'
                        icon={(
                            <Icon
                                glyph={<PlusIcon/>}
                                size='12'
                            />
                        )}
                        size='x-small'
                    />
                    <IconButton
                        aria-label='Small'
                        icon={(
                            <Icon
                                glyph={<PlusIcon/>}
                                size='16'
                            />
                        )}
                        size='small'
                    />
                    <IconButton
                        aria-label='Medium'
                        icon={icon}
                        size='medium'
                    />
                    <IconButton
                        aria-label='Large'
                        icon={(
                            <Icon
                                glyph={<PlusIcon/>}
                                size='24'
                            />
                        )}
                        size='large'
                    />
                </Row>
            </Variant>
        </>
    );
}

export function ActionButtonPreview() {
    return (
        <ActionButton
            icon={(
                <Icon
                    glyph={<CheckIcon/>}
                    size='16'
                />
            )}
            label='Message'
        />
    );
}

export function ActionButtonDetail() {
    return (
        <Variant label='Default / active / destructive'>
            <Row>
                <ActionButton
                    icon={(
                        <Icon
                            glyph={<CheckIcon/>}
                            size='16'
                        />
                    )}
                    label='Message'
                />
                <ActionButton
                    active={true}
                    icon={(
                        <Icon
                            glyph={<CheckIcon/>}
                            size='16'
                        />
                    )}
                    label='Following'
                />
                <ActionButton
                    destructive={true}
                    icon={(
                        <Icon
                            glyph={<CheckIcon/>}
                            size='16'
                        />
                    )}
                    label='Remove'
                />
            </Row>
        </Variant>
    );
}
