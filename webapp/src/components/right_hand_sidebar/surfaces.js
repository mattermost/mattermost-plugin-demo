export const sections = [
    {
        id: 'ui',
        title: 'In the UI',
        items: [
            {
                name: 'Right-hand sidebar',
                api: 'registerRightHandSidebarComponent',
                where: 'You are here.',
            },
            {
                name: 'Pop-out sidebar',
                api: 'WebappUtils.popouts.popoutRhsPlugin',
                where: 'Opens this sidebar in its own window. Not available inside a pop-out.',
                action: 'popout',
            },
            {
                name: 'App Bar button',
                api: 'registerChannelHeaderButtonAction',
                where: 'Icon in the App Bar on the far right. Opens this sidebar.',
            },
            {
                name: 'Root component',
                api: 'registerRootComponent',
                where: 'Overlay shown above everything. Most demo menu items open it.',
                action: 'openRoot',
            },
            {
                name: 'Left sidebar header',
                api: 'registerLeftSidebarHeaderComponent',
                where: 'Top of the channel sidebar, showing whether hooks are enabled.',
            },
            {
                name: 'Bottom team sidebar',
                api: 'registerBottomTeamSidebarComponent',
                where: 'Bottom of the team sidebar on the far left.',
            },
            {
                name: 'Profile popover',
                api: 'registerPopoverUserAttributesComponent, registerPopoverUserActionsComponent',
                where: 'Click any avatar or username.',
            },
            {
                name: 'Custom post types',
                api: 'registerPostTypeComponent',
                where: 'Posts in the demo channel (the ChannelName setting, demo_plugin by default).',
            },
            {
                name: 'Link tooltip',
                api: 'registerLinkTooltipComponent',
                where: 'Hover a link in any message.',
            },
            {
                name: 'Main menu',
                api: 'registerMainMenuAction',
                where: 'Click the team name at the top left: "Demo Plugin" and "Sample Confirmation Dialog".',
            },
            {
                name: 'Channel menu',
                api: 'registerChannelHeaderMenuAction',
                where: 'Click the channel name at the top: "Demo Plugin".',
            },
            {
                name: 'Message menu and submenu',
                api: 'registerPostDropdownMenuAction, registerPostDropdownSubMenuAction',
                where: 'Hover a message and open "...": "Demo Plugin" and "Submenu Example".',
            },
            {
                name: 'Browse or add channels menu',
                api: 'registerSidebarBrowseOrAddChannelMenuAction',
                where: 'The "+" above the channel list: "Demo Plugin Item".',
            },
            {
                name: 'File upload method',
                api: 'registerFileUploadMethod',
                where: 'Paperclip in the message box: "Upload using Demo Plugin".',
            },
            {
                name: 'File menu and preview',
                api: 'registerFileDropdownMenuAction, registerFilePreviewComponent',
                where: 'Upload any file ending in .demo, then open or preview it.',
            },
            {
                name: 'User settings',
                api: 'registerUserSettings',
                where: 'Settings, then the Demo Plugin tab.',
            },
            {
                name: 'System Console settings',
                api: 'registerAdminConsoleCustomSetting',
                where: 'Custom setting components on the plugin settings page.',
                action: 'openAdminConsole',
            },
        ],
    },
    {
        id: 'routes',
        title: 'Pages',
        items: [
            {
                name: 'Team page',
                api: 'registerNeedsTeamRoute',
                where: 'A plugin page inside the current team.',
                action: 'openTeamRoute',
            },
            {
                name: 'Full-screen page',
                api: 'registerCustomRoute',
                where: 'A plugin page outside the team layout.',
                action: 'openRootRoute',
            },
            {
                name: 'Sample flow page',
                api: 'registerNeedsTeamRoute',
                where: 'The sample flow\'s list and form at full width.',
                action: 'openTicketsRoute',
            },
        ],
    },
    {
        id: 'commands',
        title: 'Slash commands',
        items: [
            {name: '/demo_plugin true|false', api: 'RegisterCommand (server)', where: 'Turns the server-side hooks on or off.'},
            {name: '/dialog', api: 'RegisterCommand (server)', where: 'Opens an interactive dialog.'},
            {name: '/interactive', api: 'RegisterCommand (server)', where: 'Posts a message with interactive buttons.'},
            {name: '/ephemeral', api: 'RegisterCommand (server)', where: 'Posts a message only you can see.'},
            {name: '/ephemeral_override', api: 'RegisterCommand (server)', where: 'Ephemeral post rendered by the custom post type.'},
            {name: '/toast', api: 'RegisterCommand (server)', where: 'Shows a toast notification.'},
            {name: '/show_mentions', api: 'RegisterCommand (server)', where: 'Lists the mentions in your message.'},
            {name: '/list_files', api: 'RegisterCommand (server)', where: 'Searches files in the channel.'},
            {name: '/autocomplete_test', api: 'RegisterCommand (server)', where: 'Shows dynamic autocomplete suggestions.'},
            {name: '/crash', api: 'RegisterCommand (server)', where: 'Crashes the plugin to test restarts.'},
        ],
    },
    {
        id: 'background',
        title: 'Behind the scenes',
        items: [
            {
                name: 'Translations',
                api: 'registerTranslations',
                where: 'Strings switch with your language setting.',
                action: 'showTranslation',
            },
            {
                name: 'WebSocket events',
                api: 'registerWebSocketEventHandler',
                where: 'Updates the hooks status above when /demo_plugin is run.',
            },
            {
                name: 'Redux reducer',
                api: 'registerReducer',
                where: 'Stores the hooks status and root component state.',
            },
            {
                name: 'Reconnect handler',
                api: 'registerReconnectHandler',
                where: 'Refreshes the hooks status after a network reconnect.',
            },
        ],
    },
];

export const matchesQuery = (item, query) => {
    const q = query.trim().toLowerCase();
    if (!q) {
        return true;
    }
    return [item.name, item.api, item.where].some((text) => text.toLowerCase().includes(q));
};
