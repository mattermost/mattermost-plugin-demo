import {FormattedMessage} from 'react-intl';
import React from 'react';

export const menu = {
    id: 'primary',
    text: (
        <FormattedMessage
            id='submenu.menu'
            defaultMessage='Submenu Example'
        />
    ),
    subMenu: [
        {
            id: 'secondary.first',
            text: (
                <FormattedMessage
                    id='submenu.first'
                    key='submenu.first'
                    defaultMessage='First Item'
                />
            ),
            subMenu: [
                {
                    id: 'tertiary.first',
                    text: (
                        <FormattedMessage
                            id='submenu.second'
                            key='submenu.second'
                            defaultMessage='Second Item'
                        />
                    ),
                },
            ],
        },
        {
            id: 'secondary.second',
            text: (
                <FormattedMessage
                    id='submenu.third'
                    key='submenu.third'
                    defaultMessage='Third Item'
                />
            ),
        },
    ],
};
