import type {Store} from 'redux';

import ChannelSettingsCustomTab from './custom_tab';
import {buildSchemaTab} from './schema_tab';
import type {ChannelSettingsTab} from './types';

type Registry = {
    registerChannelSettingsTab?: (registration: ChannelSettingsTab) => void;
};

export function registerChannelSettings(registry: Registry, _store: Store) {
    if (!registry.registerChannelSettingsTab) {
        return;
    }

    registry.registerChannelSettingsTab(buildSchemaTab());

    registry.registerChannelSettingsTab({
        uiName: 'Demo Advanced (Custom)',
        icon: 'icon-tune',
        shouldRender: () => true,
        component: ChannelSettingsCustomTab,
    });
}
