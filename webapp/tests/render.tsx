// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render} from '@testing-library/react';
import type {ReactElement, ReactNode} from 'react';
import {IntlProvider} from 'react-intl';
import {MemoryRouter} from 'react-router-dom';

export function renderWithProviders(ui: ReactElement, {route = '/'} = {}) {
    return render(ui, {
        wrapper: ({children}: {children: ReactNode}) => (
            <IntlProvider locale='en'>
                <MemoryRouter initialEntries={[route]}>
                    {children}
                </MemoryRouter>
            </IntlProvider>
        ),
    });
}
