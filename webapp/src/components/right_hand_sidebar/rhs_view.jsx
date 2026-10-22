import React, {useState} from 'react';
import PropTypes from 'prop-types';
import {defineMessages, useIntl} from 'react-intl';

import {Tabs} from '@mattermost/compass-ui/components/tabs';

import {id as pluginId} from 'manifest';

import TicketFlowRhs from '../sample_flow/ticket_flow_rhs';

import ComponentsTab from './components_tab/components_tab';
import Guide from './guide';

import './rhs_view.css';

const messages = defineMessages({
    guide: {id: 'rhs.tab.guide', defaultMessage: 'Guide'},
    components: {id: 'rhs.tab.components', defaultMessage: 'Components'},
    flow: {id: 'rhs.tab.flow', defaultMessage: 'Sample flow'},
});

const TABS = [
    {key: 'guide', id: 'demo-rhs-tab-guide', panelId: 'demo-rhs-panel-guide'},
    {key: 'components', id: 'demo-rhs-tab-components', panelId: 'demo-rhs-panel-components'},
    {key: 'flow', id: 'demo-rhs-tab-flow', panelId: 'demo-rhs-panel-flow'},
];

export default function RHSView({team, channel, openRootModal}) {
    const {formatMessage} = useIntl();
    const [activeKey, setActiveKey] = useState('guide');
    const activeTab = TABS.find((tab) => tab.key === activeKey);
    const tabs = TABS.map((tab) => ({...tab, label: formatMessage(messages[tab.key])}));

    return (
        <div className='demo-rhs'>
            <Tabs
                tabs={tabs}
                activeKey={activeKey}
                onChange={setActiveKey}
                appearance='underlined'
                className='demo-rhs__tabs'
            />
            <div
                role='tabpanel'
                id={activeTab.panelId}
                aria-labelledby={activeTab.id}
                className='demo-rhs__panel'
            >
                {activeKey === 'guide' && (
                    <Guide
                        team={team}
                        channel={channel}
                        openRootModal={openRootModal}
                    />
                )}
                {activeKey === 'components' && <ComponentsTab/>}
                {activeKey === 'flow' && <TicketFlowRhs pagePath={team && `/${team.name}/${pluginId}/tickets`}/>}
            </div>
        </div>
    );
}

RHSView.propTypes = {
    team: PropTypes.object,
    channel: PropTypes.object,
    openRootModal: PropTypes.func.isRequired,
};
