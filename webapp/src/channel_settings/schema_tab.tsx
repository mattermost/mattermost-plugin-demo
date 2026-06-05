import type {Channel} from '@mattermost/types/channels';

import AccentColorSetting from './accent_color_setting';
import ChannelInfoSection from './channel_info_section';
import {fetchSchemaValues, saveSchemaValues} from './client';
import type {ChannelSettingsSchemaTab, ChannelSettingsValues} from './types';

export function buildSchemaTab(): ChannelSettingsSchemaTab {
    return {
        uiName: 'Demo Channel Settings',
        icon: 'icon-cog-outline',

        // No per-channel posting settings for DMs/GMs.
        shouldRender: (_state, channel: Channel) => channel.type !== 'D' && channel.type !== 'G',
        sections: [
            {
                title: 'Posting',
                settings: [
                    {
                        name: 'postPrefixStyle',
                        type: 'radio',
                        title: 'Post prefix style',
                        helpText: 'How the demo plugin styles its automated posts.',
                        default: 'none',
                        options: [
                            {value: 'none', text: 'None'},
                            {value: 'italic', text: 'Italic'},
                            {value: 'bold', text: 'Bold'},
                        ],
                    },
                    {
                        name: 'greetingFrequency',
                        type: 'radio',
                        title: 'Greeting frequency',
                        helpText: 'How often the demo plugin greets the channel.',
                        default: 'off',
                        options: [
                            {value: 'off', text: 'Off'},
                            {value: 'daily', text: 'Daily'},
                            {value: 'hourly', text: 'Hourly'},
                        ],
                    },
                    {
                        name: 'accentColor',
                        type: 'custom',
                        title: 'Accent color',
                        component: AccentColorSetting,
                    },
                ],
            },
            {
                title: 'Channel info',
                component: ChannelInfoSection,
            },
        ],
        loadValues: (channel: Channel) => fetchSchemaValues(channel.id),
        onSave: async (values: ChannelSettingsValues, channel: Channel) => {
            await saveSchemaValues(channel.id, values);
        },
    };
}
