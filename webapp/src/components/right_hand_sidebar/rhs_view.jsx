import React from 'react';
import PropTypes from 'prop-types';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
        unreadChannels: PropTypes.array.isRequired,
        myChannelMemberships: PropTypes.object.isRequired,
        currentUserId: PropTypes.string.isRequired,
    }

    constructor(props) {
        super(props);
        this.state = {
            contacts: [
                {id: 1, name: 'Alice Johnson', status: 'online', avatar: '👩', role: 'Developer'},
                {id: 2, name: 'Bob Smith', status: 'away', avatar: '👨', role: 'Designer'},
                {id: 3, name: 'Carol Williams', status: 'busy', avatar: '👩‍💼', role: 'Manager'},
                {id: 4, name: 'David Brown', status: 'offline', avatar: '👨‍💻', role: 'QA Engineer'},
                {id: 5, name: 'Eva Davis', status: 'online', avatar: '👩‍🔬', role: 'Scientist'},
                {id: 6, name: 'Frank Miller', status: 'away', avatar: '👨‍🎨', role: 'Artist'},
            ],
            unreadChannels: props.unreadChannels,
        };
    }

    getStatusColor = (status) => {
        switch (status) {
        case 'online':
            return '#00d100';
        case 'away':
            return '#ffaa00';
        case 'busy':
            return '#ff4444';
        case 'offline':
            return '#888888';
        default:
            return '#888888';
        }
    };

    renderContact = (contact) => {
        const statusColor = this.getStatusColor(contact.status);
        return (
            <div
                key={contact.id}
                style={styles.contactItem}
            >
                <div style={styles.contactAvatar}>
                    <span style={styles.avatarEmoji}>{contact.avatar}</span>
                    <div
                        style={{
                            ...styles.statusIndicator,
                            backgroundColor: statusColor,
                        }}
                    />
                </div>
                <div style={styles.contactInfo}>
                    <div style={styles.contactName}>{contact.name}</div>
                    <div style={styles.contactRole}>{contact.role}</div>
                </div>
                <div style={styles.contactActions}>
                    <button
                        style={styles.actionButton}
                        onClick={() => this.handleMessage(contact)}
                        title='Send message'
                    >
                        {'💬'}
                    </button>
                    <button
                        style={styles.actionButton}
                        onClick={() => this.handleCall(contact)}
                        title='Call'
                    >
                        {'📞'}
                    </button>
                </div>
            </div>
        );
    };

    handleMessage = (contact) => {
        // TODO: Implement messaging functionality
        // For now, this is a placeholder
        // Could integrate with Mattermost messaging system
        // Example: send message to contact.name
        if (contact && contact.name) {
            // Placeholder for future implementation
        }
    };

    handleCall = (contact) => {
        // TODO: Implement call functionality
        // For now, this is a placeholder
        // Could integrate with voice/video calling system
        // Example: initiate call with contact.name
        if (contact && contact.name) {
            // Placeholder for future implementation
        }
    };

    getChannelTypeLabel = (type) => {
        switch (type) {
        case 'O':
            return 'Public';
        case 'P':
            return 'Private';
        default:
            return 'Direct';
        }
    };

    handleChannelClick = (channel) => {
        // Navigate to the channel
        const teamName = this.props.team.name;
        let channelPath = '';

        // Determine the correct path based on channel type
        if (channel.type === 'D') {
            // Direct message - extract username from channel name
            // Channel name format for DMs is usually like "username1__username2"
            const usernames = channel.name.split('__');
            const otherUsername = usernames.find((username) => (
                username !== this.props.currentUserId &&
                username !== 'system'
            ));
            if (otherUsername) {
                channelPath = `/${teamName}/messages/@${otherUsername}`;
            } else {
                channelPath = `/${teamName}/messages/${channel.name}`;
            }
        } else if (channel.type === 'G') {
            // Group message
            channelPath = `/${teamName}/messages/${channel.name}`;
        } else {
            // Public/Private channel
            channelPath = `/${teamName}/channels/${channel.name}`;
        }

        // Navigate to the channel
        window.WebappUtils.browserHistory.push(channelPath);
    };

    printChannelInfo = (channel) => {
        // Print the full channel object as JSON
        // console.log('Channel Object:', JSON.stringify(channel, null, 2));

        // Show channel info in a formatted way
        const channelInfo = `
Channel Information:
==================
${JSON.stringify(channel, null, 2)}
        `.trim();

        // For now, just log to console since alert is not allowed
        // In a real implementation, you could show this in a modal or tooltip
        console.log(channelInfo); // eslint-disable-line no-console
    };

    getUnreadCount = (channel) => {
        // Get unread message count from channel membership
        const membership = this.props.myChannelMemberships[channel.id];
        if (membership) {
            return Math.max(0, channel.total_msg_count - membership.msg_count);
        }
        return 0;
    };

    componentDidUpdate(prevProps) {
        if (prevProps.unreadChannels !== this.props.unreadChannels) {
            this.setState({unreadChannels: this.props.unreadChannels});
        }
    }

    renderUnreadChannels = () => {
        const {unreadChannels} = this.state;
        if (unreadChannels.length === 0) {
            return null;
        }
        return (
            <div style={styles.unreadSection}>
                <div style={styles.unreadList}>
                    {unreadChannels.map((channel) => {
                        const unreadCount = this.getUnreadCount(channel);
                        return (
                            <div
                                key={channel.id}
                                style={{...styles.unreadItem, ...styles.unreadItemClickable}}
                                onClick={() => this.handleChannelClick(channel)}
                            >
                                <span style={styles.channelName}>
                                    {channel.display_name || channel.name}
                                </span>
                                <div style={styles.channelMeta}>
                                    <span style={styles.channelType}>
                                        {this.getChannelTypeLabel(channel.type)}
                                    </span>
                                    {unreadCount > 0 && (
                                        <span style={styles.unreadBadge}>
                                            {unreadCount > 99 ? '99+' : unreadCount}
                                        </span>
                                    )}
                                    <button
                                        style={styles.actionButton}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            this.printChannelInfo(channel);
                                        }}
                                        title='Channel Info'
                                    >
                                        {'ℹ️'}
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    };

    render() {
        const {contacts} = this.state;

        return (
            <div
                style={style.rhs}
                data-testid='rhsView'
            >

                <div style={styles.contactsList}>
                    {contacts.map(this.renderContact)}
                </div>
                {this.renderUnreadChannels()}

            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
        fontFamily: 'Open Sans, sans-serif',
    },
};

const styles = {
    title: {
        margin: '0',
        color: '#333',
        fontSize: '18px',
        fontWeight: '600',
    },
    contactsList: {
        marginBottom: '20px',
    },
    contactItem: {
        display: 'flex',
        alignItems: 'center',
        padding: '10px 0',
        borderBottom: '1px solid #f0f0f0',
        transition: 'background-color 0.2s ease',
        cursor: 'pointer',
        ':hover': {
            backgroundColor: '#f8f8f8',
        },
    },
    contactAvatar: {
        position: 'relative',
        marginRight: '12px',
    },
    avatarEmoji: {
        fontSize: '32px',
        display: 'block',
    },
    statusIndicator: {
        position: 'absolute',
        bottom: '2px',
        right: '2px',
        width: '12px',
        height: '12px',
        borderRadius: '50%',
        border: '2px solid white',
    },
    contactInfo: {
        flex: 1,
    },
    contactName: {
        fontWeight: '500',
        color: '#333',
        marginBottom: '2px',
    },
    contactRole: {
        fontSize: '12px',
        color: '#666',
    },
    contactActions: {
        display: 'flex',
        gap: '5px',
    },
    actionButton: {
        background: 'none',
        border: 'none',
        cursor: 'pointer',
        padding: '5px',
        borderRadius: '3px',
        fontSize: '14px',
        transition: 'background-color 0.2s ease',
        ':hover': {
            backgroundColor: '#e6f3ff',
        },
    },
    footer: {
        borderTop: '1px solid #e1e1e1',
        paddingTop: '10px',
        marginBottom: '15px',
    },
    footerText: {
        margin: '0',
        fontSize: '12px',
        color: '#666',
        textAlign: 'center',
    },
    links: {
        borderTop: '1px solid #e1e1e1',
        paddingTop: '15px',
    },
    linksTitle: {
        margin: '0 0 10px 0',
        fontSize: '14px',
        color: '#666',
        fontWeight: '500',
    },
    link: {
        color: '#0066cc',
        textDecoration: 'none',
        fontSize: '12px',
        cursor: 'pointer',
        ':hover': {
            textDecoration: 'underline',
        },
    },
    unreadSection: {
        marginTop: '20px',
    },
    unreadTitle: {
        margin: '0',
        fontSize: '16px',
        color: '#333',
        fontWeight: '500',
    },
    unreadList: {
        listStyle: 'none',
        padding: '0',
        margin: '0',
    },
    unreadItem: {
        padding: '10px',
        borderBottom: '1px solid #f0f0f0',
    },
    unreadItemClickable: {
        cursor: 'pointer',
        transition: 'background-color 0.2s ease',
        ':hover': {
            backgroundColor: '#f0f8ff',
        },
    },
    channelName: {
        fontSize: '14px',
        color: '#333',
    },
    channelType: {
        fontSize: '12px',
        color: '#666',
    },
    unreadCountBadge: {
        fontSize: '12px',
        color: '#666',
        marginLeft: '5px',
    },
    channelMeta: {
        fontSize: '12px',
        color: '#666',
        display: 'flex',
        alignItems: 'center',
        gap: '5px',
    },
    unreadBadge: {
        fontSize: '12px',
        color: '#fff',
        backgroundColor: '#ff4444',
        padding: '2px 4px',
        borderRadius: '3px',
        marginLeft: '5px',
    },
};
