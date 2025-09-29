import React from 'react';

import {FormattedMessage} from 'react-intl';

import en from 'i18n/en.json';

import es from 'i18n/es.json';

import manifest from './manifest';

import Root from './components/root';
import BottomTeamSidebar from './components/bottom_team_sidebar';
import LeftSidebarHeader from './components/left_sidebar_header';
import LinkTooltip from './components/link_tooltip';
import UserAttributes from './components/user_attributes';
import UserActions from './components/user_actions';
import RHSView from './components/right_hand_sidebar';
import SecretMessageSetting from './components/admin_settings/secret_message_setting';
import CustomSetting from './components/admin_settings/custom_setting';
import FilePreviewOverride from './components/file_preview_override';
import RouterShowcase from './components/router_showcase/router_showcase';
import PostType from './components/post_type';
import EphemeralPostType from './components/ephemeral_post_type';
import {
    MainMenuMobileIcon,
    ChannelHeaderButtonIcon,
    FileUploadMethodIcon,
} from './components/icons';
import {
    mainMenuAction,
    fileDropdownMenuAction,
    fileUploadMethodAction,
    postDropdownMenuAction,
    postDropdownSubMenuAction,
    channelHeaderMenuAction,
    websocketStatusChange,
    getStatus,
} from './actions';
import reducer from './reducer';

function getTranslations(locale) {
    switch (locale) {
    case 'en':
        return en;
    case 'es':
        return es;
    }
    return {};
}

export default class DemoPlugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);
        registry.registerPopoverUserAttributesComponent(UserAttributes);
        registry.registerPopoverUserActionsComponent(UserActions);
        registry.registerLeftSidebarHeaderComponent(LeftSidebarHeader);
        registry.registerLinkTooltipComponent(LinkTooltip);
        registry.registerBottomTeamSidebarComponent(
            BottomTeamSidebar,
        );
        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            RHSView,
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WhatsApp'
            />);

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(toggleRHSPlugin),
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WhatsApp'
            />,
        );

        registry.registerPostTypeComponent('custom_demo_plugin', PostType);
        registry.registerPostTypeComponent('custom_demo_plugin_ephemeral', EphemeralPostType);

        registry.registerMainMenuAction(
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WhatsApp'
            />,
            () => store.dispatch(mainMenuAction()),
            <MainMenuMobileIcon/>,
        );

        registry.registerChannelHeaderMenuAction(
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WhatsApp'
            />,
            (channelId) => store.dispatch(channelHeaderMenuAction(channelId)),
            <MainMenuMobileIcon/>,
        );

        registry.registerMainMenuAction(
            <FormattedMessage
                id='sample.confirmation.dialog'
                defaultMessage='Sample Confirmation Dialog'
            />,
            () => {
                window.openInteractiveDialog({
                    url: '/plugins/' + manifest.id + '/dialog/2',
                    dialog: {
                        callback_id: 'somecallbackid',
                        title: 'Sample Confirmation Dialog',
                        elements: [],
                        submit_label: 'Confirm',
                        notify_on_cancel: true,
                        state: 'somestate',
                    },
                });
            },
            <MainMenuMobileIcon/>,
        );

        registry.registerPostDropdownMenuAction(
            <FormattedMessage
                id='plugin.name'
                defaultMessage='WhatsApp'
            />,
            () => store.dispatch(postDropdownMenuAction()),
        );

        // eslint-disable-next-line no-unused-vars
        const {id, rootRegisterMenuItem} = registry.registerPostDropdownSubMenuAction(
            <FormattedMessage
                id='submenu.menu'
                key='submenu.menu'
                defaultMessage='Submenu Example'
            />,
        );

        const firstItem = (
            <FormattedMessage
                id='submenu.first'
                key='submenu.first'
                defaultMessage='First Item'
            />
        );
        rootRegisterMenuItem(
            firstItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(firstItem));
            },
        );

        const secondItem = (
            <FormattedMessage
                id='submenu.second'
                key='submenu.second'
                defaultMessage='Second Item'
            />
        );
        rootRegisterMenuItem(
            secondItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(secondItem));
            },
        );

        const thirdItem = (
            <FormattedMessage
                id='submenu.third'
                key='submenu.third'
                defaultMessage='Third Item'
            />
        );
        rootRegisterMenuItem(
            thirdItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(thirdItem));
            },
        );

        registry.registerFileUploadMethod(
            <FileUploadMethodIcon/>,
            () => store.dispatch(fileUploadMethodAction()),
            <FormattedMessage
                id='plugin.upload'
                defaultMessage='WhatsApp: Notificar package'
            />,
        );

        // ignore if registerFileDropdownMenuAction method does not exist
        if (registry.registerFileDropdownMenuAction) {
            registry.registerFileDropdownMenuAction(
                (fileInfo) => fileInfo.extension === 'demo',
                <FormattedMessage
                    id='plugin.name'
                    defaultMessage='WhatsApp'
                />,
                () => store.dispatch(fileDropdownMenuAction()),
            );
        }

        registry.registerWebSocketEventHandler(
            'custom_' + manifest.id + '_status_change',
            (message) => {
                store.dispatch(websocketStatusChange(message));
            },
        );

        registry.registerAdminConsoleCustomSetting('SecretMessage', SecretMessageSetting, {showTitle: true});
        registry.registerAdminConsoleCustomSetting('CustomSetting', CustomSetting);

        registry.registerFilePreviewComponent((fileInfo) => fileInfo.extension === 'demo', FilePreviewOverride);

        registry.registerReducer(reducer);

        // Immediately fetch the current plugin status.
        store.dispatch(getStatus());

        // Fetch the current status whenever we recover an internet connection.
        registry.registerReconnectHandler(() => {
            store.dispatch(getStatus());
        });

        registry.registerTranslations(getTranslations);

        registry.registerNeedsTeamRoute('/teamtest', RouterShowcase);
        registry.registerCustomRoute('/roottest', () => 'Demo plugin route.');

        registry.registerUserSettings?.({
            id: manifest.id,
            icon: `/plugins/${manifest.id}/public/icon.png`,
            uiName: manifest.name,
            action: {
                title: 'Example action',
                text: 'This is an example action for this setting',
                buttonText: 'Here is the button text',
                onClick: () => alert('Button clicked'), // eslint-disable-line no-alert
            },
            sections: [
                {
                    settings: [
                        {
                            name: 'setting1',
                            title: 'setting 1',
                            options: [
                                {
                                    text: 'Option 1',
                                    value: '1',
                                    helpText: 'Some option help text',
                                },
                                {
                                    text: 'Option 2',
                                    value: '2',
                                },
                            ],
                            type: 'radio',
                            default: '1',
                            helpText: 'Some setting help text',
                        },
                        {
                            name: 'setting2',
                            title: 'setting 2',
                            options: [
                                {
                                    text: 'Option 1',
                                    value: '1',
                                    helpText: 'Some option help text',
                                },
                                {
                                    text: 'Option 2',
                                    value: '2',
                                },
                                {
                                    text: 'Option 3',
                                    value: '3',
                                    helpText: 'Some option help text',
                                },
                            ],
                            type: 'radio',
                            default: '1',
                        },
                    ],
                    title: 'Test section number 1',
                    onSubmit: (v) => alert(`saving ${Object.keys(v).map((k) => `{${k}}: ${v[k]}`).join(' ')}`), // eslint-disable-line no-alert
                },
                {
                    settings: [{
                        name: 'setting3',
                        options: [
                            {
                                text: 'Option 1',
                                value: '1',
                                helpText: 'Some option help text',
                            },
                            {
                                text: 'Option 2',
                                value: '2',
                            },
                        ],
                        type: 'radio',
                        default: '2',
                    }],
                    title: 'Test section number 2',
                    onSubmit: (v) => alert(`saving ${Object.keys(v).map((k) => `{${k}}: ${v[k]}`).join(' ')}`), // eslint-disable-line no-alert
                },
                {
                    settings: [{
                        name: 'setting4',
                        options: [
                            {
                                text: 'Option 1',
                                value: '1',
                                helpText: 'Some option help text',
                            },
                            {
                                text: 'Option 2',
                                value: '2',
                            },
                        ],
                        type: 'radio',
                        default: '2',
                    }],
                    title: 'Test section disabled',
                    disabled: true,
                    onSubmit: (v) => alert(`saving ${Object.keys(v).map((k) => `{${k}}: ${v[k]}`).join(' ')}`), // eslint-disable-line no-alert
                },
            ],
        });
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(manifest.id + '::uninitialize()');
    }
}
