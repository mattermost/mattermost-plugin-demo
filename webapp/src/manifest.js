// This file is automatically generated. Do not modify it manually.

const manifest = JSON.parse(`
{
    "id": "com.mattermost.demo-plugin",
    "name": "Demo Plugin",
    "description": "This plugin demonstrates the capabilities of a Mattermost plugin.",
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-demo",
    "support_url": "https://github.com/mattermost/mattermost-plugin-demo/issues",
    "icon_path": "assets/icon.svg",
    "version": "0.8.0",
    "min_server_version": "5.26.0",
    "server": {
        "executables": {
            "linux-amd64": "server/dist/plugin-linux-amd64",
            "darwin-amd64": "server/dist/plugin-darwin-amd64",
            "windows-amd64": "server/dist/plugin-windows-amd64.exe"
        },
        "executable": ""
    },
    "webapp": {
        "bundle_path": "webapp/dist/main.js"
    },
    "settings_schema": {
        "header": "Header: Configure your demo plugin settings below",
        "footer": "Footer: The code for this demo plugin can be found [here](https://github.com/mattermost/mattermost-plugin-demo)",
        "settings": [
            {
                "key": "ChannelName",
                "display_name": "Channel Name",
                "type": "text",
                "help_text": "The channel to use as part of the demo plugin, created for each team automatically if it does not exist.",
                "placeholder": "demo_plugin",
                "default": "demo_plugin"
            },
            {
                "key": "Username",
                "display_name": "Username",
                "type": "text",
                "help_text": "The user to use as part of the demo plugin, created automatically if it does not exist.",
                "placeholder": "demo_plugin",
                "default": "demo_plugin"
            },
            {
                "key": "LastName",
                "display_name": "Demo User Last Name",
                "type": "radio",
                "help_text": "Select the last name for the demo user",
                "placeholder": "",
                "default": "Plugin User",
                "options": [
                    {
                        "display_name": "Plugin User",
                        "value": "Plugin User"
                    },
                    {
                        "display_name": "Demoson III",
                        "value": "Demoson III"
                    },
                    {
                        "display_name": "McDemo",
                        "value": "McDemo"
                    }
                ]
            },
            {
                "key": "TextStyle",
                "display_name": "Text Style",
                "type": "dropdown",
                "help_text": "Change the text style of the messages posted by this plugin",
                "placeholder": "",
                "default": "",
                "options": [
                    {
                        "display_name": "none",
                        "value": ""
                    },
                    {
                        "display_name": "italics",
                        "value": "_"
                    },
                    {
                        "display_name": "bold",
                        "value": "**"
                    }
                ]
            },
            {
                "key": "RandomSecret",
                "display_name": "Random Secret",
                "type": "generated",
                "help_text": "Generate a random string that the demo plugin will watch for. If the secret string is mentioned in any channel then the demo plugin will publish a special message.",
                "regenerate_help_text": "Generate a new secret string.",
                "placeholder": "",
                "default": "CFgcq9Hr9OKSevvqH_SH-mPlgVklmpUm"
            },
            {
                "key": "SecretMessage",
                "display_name": "Secret Message",
                "type": "custom",
                "help_text": "The message posted by the demo plugin when the secret phrase is detected.",
                "placeholder": "",
                "default": "Yay! The random secret string was posted! Go to the settings page for this plugin in the system console to generate a new random secret."
            },
            {
                "key": "CustomSetting",
                "display_name": "",
                "type": "custom",
                "help_text": "",
                "placeholder": "",
                "default": null
            },
            {
                "key": "EnableMentionUser",
                "display_name": "Enable Mention User",
                "type": "bool",
                "help_text": "Enable or disable the demo plugin to tag a username on every message sent. The username value is set below.",
                "placeholder": "",
                "default": false
            },
            {
                "key": "MentionUser",
                "display_name": "Mention User",
                "type": "username",
                "help_text": "Configure the username to be mentioned by the demo plugin. Must be enabled in the setting above.",
                "placeholder": "demo_plugin",
                "default": "demo_plugin"
            },
            {
                "key": "secretNumber",
                "display_name": "Secret Number",
                "type": "number",
                "help_text": "A secret number that the demo plugin will watch for. If the secret number is mentioned in any channel then the demo plugin will publish a special message.",
                "placeholder": "Some secret number",
                "default": 123
            }
        ]
    },
    "required_configuration": {
        "ServiceSettings": {
            "SiteURL": null,
            "WebsocketURL": null,
            "LicenseFileLocation": null,
            "ListenAddress": null,
            "ConnectionSecurity": null,
            "TLSCertFile": null,
            "TLSKeyFile": null,
            "TLSMinVer": null,
            "TLSStrictTransport": null,
            "TLSStrictTransportMaxAge": null,
            "TLSOverwriteCiphers": null,
            "UseLetsEncrypt": null,
            "LetsEncryptCertificateCacheFile": null,
            "Forward80To443": null,
            "TrustedProxyIPHeader": null,
            "ReadTimeout": null,
            "WriteTimeout": null,
            "IdleTimeout": null,
            "MaximumLoginAttempts": null,
            "GoroutineHealthThreshold": null,
            "GoogleDeveloperKey": null,
            "EnableOAuthServiceProvider": null,
            "EnableIncomingWebhooks": null,
            "EnableOutgoingWebhooks": null,
            "EnableCommands": null,
            "EnableOnlyAdminIntegrations": null,
            "EnablePostUsernameOverride": null,
            "EnablePostIconOverride": null,
            "EnableLinkPreviews": null,
            "EnableTesting": null,
            "EnableDeveloper": null,
            "EnableOpenTracing": null,
            "EnableSecurityFixAlert": null,
            "EnableInsecureOutgoingConnections": null,
            "AllowedUntrustedInternalConnections": null,
            "EnableMultifactorAuthentication": null,
            "EnforceMultifactorAuthentication": null,
            "EnableUserAccessTokens": null,
            "AllowCorsFrom": null,
            "CorsExposedHeaders": null,
            "CorsAllowCredentials": null,
            "CorsDebug": null,
            "AllowCookiesForSubdomains": null,
            "ExtendSessionLengthWithActivity": null,
            "SessionLengthWebInDays": null,
            "SessionLengthMobileInDays": null,
            "SessionLengthSSOInDays": null,
            "SessionCacheInMinutes": null,
            "SessionIdleTimeoutInMinutes": null,
            "WebsocketSecurePort": null,
            "WebsocketPort": null,
            "WebserverMode": null,
            "EnableCustomEmoji": null,
            "EnableEmojiPicker": null,
            "EnableGifPicker": true,
            "GfycatApiKey": null,
            "GfycatApiSecret": null,
            "RestrictCustomEmojiCreation": null,
            "RestrictPostDelete": null,
            "AllowEditPost": null,
            "PostEditTimeLimit": null,
            "TimeBetweenUserTypingUpdatesMilliseconds": null,
            "EnablePostSearch": null,
            "MinimumHashtagLength": null,
            "EnableUserTypingMessages": null,
            "EnableChannelViewedMessages": null,
            "EnableUserStatuses": null,
            "ExperimentalEnableAuthenticationTransfer": null,
            "ClusterLogTimeoutMilliseconds": null,
            "CloseUnusedDirectMessages": null,
            "EnablePreviewFeatures": null,
            "EnableTutorial": null,
            "ExperimentalEnableDefaultChannelLeaveJoinMessages": null,
            "ExperimentalGroupUnreadChannels": null,
            "ExperimentalChannelOrganization": null,
            "ExperimentalChannelSidebarOrganization": null,
            "ExperimentalDataPrefetch": null,
            "ImageProxyType": null,
            "ImageProxyURL": null,
            "ImageProxyOptions": null,
            "EnableAPITeamDeletion": null,
            "EnableAPIUserDeletion": null,
            "ExperimentalEnableHardenedMode": null,
            "DisableLegacyMFA": null,
            "ExperimentalStrictCSRFEnforcement": null,
            "EnableEmailInvitations": null,
            "DisableBotsWhenOwnerIsDeactivated": null,
            "EnableBotAccountCreation": null,
            "EnableSVGs": null,
            "EnableLatex": null,
            "EnableAPIChannelDeletion": null,
            "EnableLocalMode": null,
            "LocalModeSocketLocation": null,
            "EnableAWSMetering": null
        },
        "TeamSettings": {
            "SiteName": null,
            "MaxUsersPerTeam": null,
            "EnableTeamCreation": null,
            "EnableUserCreation": null,
            "EnableOpenServer": null,
            "EnableUserDeactivation": null,
            "RestrictCreationToDomains": null,
            "EnableCustomBrand": null,
            "CustomBrandText": null,
            "CustomDescriptionText": null,
            "RestrictDirectMessage": null,
            "RestrictTeamInvite": null,
            "RestrictPublicChannelManagement": null,
            "RestrictPrivateChannelManagement": null,
            "RestrictPublicChannelCreation": null,
            "RestrictPrivateChannelCreation": null,
            "RestrictPublicChannelDeletion": null,
            "RestrictPrivateChannelDeletion": null,
            "RestrictPrivateChannelManageMembers": null,
            "EnableXToLeaveChannelsFromLHS": null,
            "UserStatusAwayTimeout": null,
            "MaxChannelsPerTeam": null,
            "MaxNotificationsPerChannel": null,
            "EnableConfirmNotificationsToChannel": null,
            "TeammateNameDisplay": null,
            "ExperimentalViewArchivedChannels": null,
            "ExperimentalEnableAutomaticReplies": null,
            "ExperimentalHideTownSquareinLHS": null,
            "ExperimentalTownSquareIsReadOnly": null,
            "LockTeammateNameDisplay": null,
            "ExperimentalPrimaryTeam": null,
            "ExperimentalDefaultChannels": null
        },
        "ClientRequirements": {
            "AndroidLatestVersion": "",
            "AndroidMinVersion": "",
            "DesktopLatestVersion": "",
            "DesktopMinVersion": "",
            "IosLatestVersion": "",
            "IosMinVersion": ""
        },
        "SqlSettings": {
            "DriverName": null,
            "DataSource": null,
            "DataSourceReplicas": null,
            "DataSourceSearchReplicas": null,
            "MaxIdleConns": null,
            "ConnMaxLifetimeMilliseconds": null,
            "MaxOpenConns": null,
            "Trace": null,
            "AtRestEncryptKey": null,
            "QueryTimeout": null,
            "DisableDatabaseSearch": null
        },
        "LogSettings": {
            "EnableConsole": null,
            "ConsoleLevel": null,
            "ConsoleJson": null,
            "EnableFile": null,
            "FileLevel": null,
            "FileJson": null,
            "FileLocation": null,
            "EnableWebhookDebugging": null,
            "EnableDiagnostics": null,
            "EnableSentry": null,
            "AdvancedLoggingConfig": null
        },
        "ExperimentalAuditSettings": {
            "FileEnabled": null,
            "FileName": null,
            "FileMaxSizeMB": null,
            "FileMaxAgeDays": null,
            "FileMaxBackups": null,
            "FileCompress": null,
            "FileMaxQueueSize": null,
            "AdvancedLoggingConfig": null
        },
        "NotificationLogSettings": {
            "EnableConsole": null,
            "ConsoleLevel": null,
            "ConsoleJson": null,
            "EnableFile": null,
            "FileLevel": null,
            "FileJson": null,
            "FileLocation": null,
            "AdvancedLoggingConfig": null
        },
        "PasswordSettings": {
            "MinimumLength": null,
            "Lowercase": null,
            "Number": null,
            "Uppercase": null,
            "Symbol": null
        },
        "FileSettings": {
            "EnableFileAttachments": null,
            "EnableMobileUpload": null,
            "EnableMobileDownload": null,
            "MaxFileSize": null,
            "DriverName": null,
            "Directory": null,
            "EnablePublicLink": true,
            "PublicLinkSalt": null,
            "InitialFont": null,
            "AmazonS3AccessKeyId": null,
            "AmazonS3SecretAccessKey": null,
            "AmazonS3Bucket": null,
            "AmazonS3PathPrefix": null,
            "AmazonS3Region": null,
            "AmazonS3Endpoint": null,
            "AmazonS3SSL": null,
            "AmazonS3SignV2": null,
            "AmazonS3SSE": null,
            "AmazonS3Trace": null
        },
        "EmailSettings": {
            "EnableSignUpWithEmail": null,
            "EnableSignInWithEmail": null,
            "EnableSignInWithUsername": null,
            "SendEmailNotifications": null,
            "UseChannelInEmailNotifications": null,
            "RequireEmailVerification": null,
            "FeedbackName": null,
            "FeedbackEmail": null,
            "ReplyToAddress": null,
            "FeedbackOrganization": null,
            "EnableSMTPAuth": null,
            "SMTPUsername": null,
            "SMTPPassword": null,
            "SMTPServer": null,
            "SMTPPort": null,
            "SMTPServerTimeout": null,
            "ConnectionSecurity": null,
            "SendPushNotifications": null,
            "PushNotificationServer": null,
            "PushNotificationContents": null,
            "PushNotificationBuffer": null,
            "EnableEmailBatching": null,
            "EmailBatchingBufferSize": null,
            "EmailBatchingInterval": null,
            "EnablePreviewModeBanner": null,
            "SkipServerCertificateVerification": null,
            "EmailNotificationContentsType": null,
            "LoginButtonColor": null,
            "LoginButtonBorderColor": null,
            "LoginButtonTextColor": null
        },
        "RateLimitSettings": {
            "Enable": null,
            "PerSec": null,
            "MaxBurst": null,
            "MemoryStoreSize": null,
            "VaryByRemoteAddr": null,
            "VaryByUser": null,
            "VaryByHeader": ""
        },
        "PrivacySettings": {
            "ShowEmailAddress": null,
            "ShowFullName": null
        },
        "SupportSettings": {
            "TermsOfServiceLink": null,
            "PrivacyPolicyLink": null,
            "AboutLink": null,
            "HelpLink": null,
            "ReportAProblemLink": null,
            "SupportEmail": null,
            "CustomTermsOfServiceEnabled": null,
            "CustomTermsOfServiceReAcceptancePeriod": null,
            "EnableAskCommunityLink": null
        },
        "AnnouncementSettings": {
            "EnableBanner": null,
            "BannerText": null,
            "BannerColor": null,
            "BannerTextColor": null,
            "AllowBannerDismissal": null,
            "AdminNoticesEnabled": null,
            "UserNoticesEnabled": null,
            "NoticesURL": null,
            "NoticesFetchFrequency": null,
            "NoticesSkipCache": null
        },
        "ThemeSettings": {
            "EnableThemeSelection": null,
            "DefaultTheme": null,
            "AllowCustomThemes": null,
            "AllowedThemes": null
        },
        "GitLabSettings": {
            "Enable": null,
            "Secret": null,
            "Id": null,
            "Scope": null,
            "AuthEndpoint": null,
            "TokenEndpoint": null,
            "UserApiEndpoint": null
        },
        "GoogleSettings": {
            "Enable": null,
            "Secret": null,
            "Id": null,
            "Scope": null,
            "AuthEndpoint": null,
            "TokenEndpoint": null,
            "UserApiEndpoint": null
        },
        "Office365Settings": {
            "Enable": null,
            "Secret": null,
            "Id": null,
            "Scope": null,
            "AuthEndpoint": null,
            "TokenEndpoint": null,
            "UserApiEndpoint": null,
            "DirectoryId": null
        },
        "LdapSettings": {
            "Enable": null,
            "EnableSync": null,
            "LdapServer": null,
            "LdapPort": null,
            "ConnectionSecurity": null,
            "BaseDN": null,
            "BindUsername": null,
            "BindPassword": null,
            "UserFilter": null,
            "GroupFilter": null,
            "GuestFilter": null,
            "EnableAdminFilter": null,
            "AdminFilter": null,
            "GroupDisplayNameAttribute": null,
            "GroupIdAttribute": null,
            "FirstNameAttribute": null,
            "LastNameAttribute": null,
            "EmailAttribute": null,
            "UsernameAttribute": null,
            "NicknameAttribute": null,
            "IdAttribute": null,
            "PositionAttribute": null,
            "LoginIdAttribute": null,
            "PictureAttribute": null,
            "SyncIntervalMinutes": null,
            "SkipCertificateVerification": null,
            "PublicCertificateFile": null,
            "PrivateKeyFile": null,
            "QueryTimeout": null,
            "MaxPageSize": null,
            "LoginFieldName": null,
            "LoginButtonColor": null,
            "LoginButtonBorderColor": null,
            "LoginButtonTextColor": null,
            "Trace": null
        },
        "ComplianceSettings": {
            "Enable": null,
            "Directory": null,
            "EnableDaily": null
        },
        "LocalizationSettings": {
            "DefaultServerLocale": null,
            "DefaultClientLocale": null,
            "AvailableLocales": null
        },
        "SamlSettings": {
            "Enable": null,
            "EnableSyncWithLdap": null,
            "EnableSyncWithLdapIncludeAuth": null,
            "Verify": null,
            "Encrypt": null,
            "SignRequest": null,
            "IdpUrl": null,
            "IdpDescriptorUrl": null,
            "IdpMetadataUrl": null,
            "ServiceProviderIdentifier": null,
            "AssertionConsumerServiceURL": null,
            "SignatureAlgorithm": null,
            "CanonicalAlgorithm": null,
            "ScopingIDPProviderId": null,
            "ScopingIDPName": null,
            "IdpCertificateFile": null,
            "PublicCertificateFile": null,
            "PrivateKeyFile": null,
            "IdAttribute": null,
            "GuestAttribute": null,
            "EnableAdminAttribute": null,
            "AdminAttribute": null,
            "FirstNameAttribute": null,
            "LastNameAttribute": null,
            "EmailAttribute": null,
            "UsernameAttribute": null,
            "NicknameAttribute": null,
            "LocaleAttribute": null,
            "PositionAttribute": null,
            "LoginButtonText": null,
            "LoginButtonColor": null,
            "LoginButtonBorderColor": null,
            "LoginButtonTextColor": null
        },
        "NativeAppSettings": {
            "AppDownloadLink": null,
            "AndroidAppDownloadLink": null,
            "IosAppDownloadLink": null
        },
        "ClusterSettings": {
            "Enable": null,
            "ClusterName": null,
            "OverrideHostname": null,
            "NetworkInterface": null,
            "BindAddress": null,
            "AdvertiseAddress": null,
            "UseIpAddress": null,
            "UseExperimentalGossip": null,
            "EnableExperimentalGossipEncryption": null,
            "ReadOnlyConfig": null,
            "GossipPort": null,
            "StreamingPort": null,
            "MaxIdleConns": null,
            "MaxIdleConnsPerHost": null,
            "IdleConnTimeoutMilliseconds": null
        },
        "MetricsSettings": {
            "Enable": null,
            "BlockProfileRate": null,
            "ListenAddress": null
        },
        "ExperimentalSettings": {
            "ClientSideCertEnable": null,
            "ClientSideCertCheck": null,
            "EnableClickToReply": null,
            "LinkMetadataTimeoutMilliseconds": null,
            "RestrictSystemAdmin": null,
            "UseNewSAMLLibrary": null,
            "CloudUserLimit": null,
            "CloudBilling": null,
            "EnableSharedChannels": null
        },
        "AnalyticsSettings": {
            "MaxUsersForStatistics": null
        },
        "ElasticsearchSettings": {
            "ConnectionUrl": null,
            "Username": null,
            "Password": null,
            "EnableIndexing": null,
            "EnableSearching": null,
            "EnableAutocomplete": null,
            "Sniff": null,
            "PostIndexReplicas": null,
            "PostIndexShards": null,
            "ChannelIndexReplicas": null,
            "ChannelIndexShards": null,
            "UserIndexReplicas": null,
            "UserIndexShards": null,
            "AggregatePostsAfterDays": null,
            "PostsAggregatorJobStartTime": null,
            "IndexPrefix": null,
            "LiveIndexingBatchSize": null,
            "BulkIndexingTimeWindowSeconds": null,
            "RequestTimeoutSeconds": null,
            "SkipTLSVerification": null,
            "Trace": null
        },
        "BleveSettings": {
            "IndexDir": null,
            "EnableIndexing": null,
            "EnableSearching": null,
            "EnableAutocomplete": null,
            "BulkIndexingTimeWindowSeconds": null
        },
        "DataRetentionSettings": {
            "EnableMessageDeletion": null,
            "EnableFileDeletion": null,
            "MessageRetentionDays": null,
            "FileRetentionDays": null,
            "DeletionJobStartTime": null
        },
        "MessageExportSettings": {
            "EnableExport": null,
            "ExportFormat": null,
            "DailyRunTime": null,
            "ExportFromTimestamp": null,
            "BatchSize": null,
            "DownloadExportResults": null,
            "GlobalRelaySettings": null
        },
        "JobSettings": {
            "RunJobs": null,
            "RunScheduler": null
        },
        "PluginSettings": {
            "Enable": null,
            "EnableUploads": null,
            "AllowInsecureDownloadUrl": null,
            "EnableHealthCheck": null,
            "Directory": null,
            "ClientDirectory": null,
            "Plugins": null,
            "PluginStates": null,
            "EnableMarketplace": null,
            "EnableRemoteMarketplace": null,
            "AutomaticPrepackagedPlugins": null,
            "RequirePluginSignature": null,
            "MarketplaceUrl": null,
            "SignaturePublicKeyFiles": null
        },
        "DisplaySettings": {
            "CustomUrlSchemes": null,
            "ExperimentalTimezone": null
        },
        "GuestAccountsSettings": {
            "Enable": null,
            "AllowEmailAccounts": null,
            "EnforceMultifactorAuthentication": null,
            "RestrictCreationToDomains": null
        },
        "ImageProxySettings": {
            "Enable": null,
            "ImageProxyType": null,
            "RemoteImageProxyURL": null,
            "RemoteImageProxyOptions": null
        },
        "CloudSettings": {
            "CWSUrl": null
        }
    }
}
`);

export default manifest;
export const id = manifest.id;
export const version = manifest.version;
