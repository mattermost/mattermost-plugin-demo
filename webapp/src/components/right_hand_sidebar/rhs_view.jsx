import React from 'react';
import PropTypes from 'prop-types';
import {getStatusForUserId} from 'mattermost-redux/selectors/entities/users';

import {getBasePath} from '../../actions';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
        unreadChannels: PropTypes.array.isRequired,
        myChannelMemberships: PropTypes.object.isRequired,
        currentUserId: PropTypes.string.isRequired,
        activeUsers: PropTypes.array.isRequired,
        reduxState: PropTypes.object.isRequired,
    }

    constructor(props) {
        super(props);
        this.state = {
            unreadChannels: props.unreadChannels,
        };
    }

    getStatusStyle = (status) => {
        switch (status) {
        case 'online':
            return 'sc-dkrFOg jcrjOR sc-kgTSHT dKNYAk icon-check-circle';
        case 'away':
            return 'sc-dkrFOg jcrjOR sc-kgTSHT iEyIcE icon-clock';
        case 'dnd':
            return 'sc-dkrFOg jcrjOR sc-kgTSHT fFlKgj icon-minus-circle';
        case 'offline':
            return 'sc-dkrFOg jcrjOR sc-kgTSHT hRRwBM icon-circle-outline';
        default:
            return '#888888';
        }
    };

    renderContact = (contact) => {
        const state = this.props.reduxState;
        const userStatus = getStatusForUserId(state, contact.id);
        const statusStyle = this.getStatusStyle(userStatus);
        const defaultStatusStyle = {
            border: '8px solid white',
            borderRadius: '50%',
        };
        return (
            <div
                key={contact.id}
                style={styles.contactCard}
            >
                <div
                    style={styles.avatarCircle}
                    title={contact.username}
                    aria-label={contact.username}
                >
                    <a href={`${getBasePath(this.props.reduxState)}/${this.props.team?.name}/messages/@${contact.username}`}>
                        <img
                            src={`${getBasePath(this.props.reduxState)}/api/v4/users/${contact.id}/image?_=0`}
                            style={styles.avatar}
                        />
                    </a>
                    <div style={styles.ballStatusStyle}>
                        <i
                            className={statusStyle}
                            style={defaultStatusStyle}
                        />
                    </div>
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
                        const statusColor = '#ff4444';
                        return (
                            <div
                                key={channel.id}
                                style={styles.contactItem}
                                onClick={() => this.handleChannelClick(channel)}
                            >
                                <div style={styles.contactAvatar}>
                                    <span style={styles.avatarEmoji}>{'👩'}</span>
                                    <div
                                        style={{
                                            ...styles.statusIndicator,
                                            backgroundColor: statusColor,
                                        }}
                                    />
                                </div>
                                <div style={styles.contactInfo}>
                                    <div style={styles.contactName}>
                                        {channel.display_name || channel.name}
                                    </div>
                                    <div style={styles.contactRole}>
                                        {'New message'}
                                        {unreadCount > 0 && (
                                            <span style={styles.unreadBadge}>
                                                {unreadCount > 99 ? '99+' : unreadCount}
                                            </span>
                                        )}
                                    </div>
                                </div>
                                <div style={styles.contactActions}>
                                    <button
                                        style={styles.actionButton}
                                        onClick={() => this.handleChannelClick(channel)}
                                        title='Send message'
                                    >
                                        {'👀'}
                                    </button>
                                    <button
                                        style={styles.actionButton}
                                        onClick={() => {
                                            this.handleCall(channel);
                                        }}
                                        title='Call'
                                    >
                                        {'✅'}
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
        const {activeUsers} = this.props;

        return (
            <div
                style={style.rhs}
                data-testid='rhsView'
            >
                <h3 style={styles.viewTitle}>{'Atención al cliente'}</h3>
                <div style={{...styles.contactsList, maxHeight: '300px', overflowY: 'auto', borderBottom: '1px solid #dddddd'}}>
                    {activeUsers?.map(this.renderContact)}
                </div>
                <h4 style={styles.chatsTitle}>{'Chats'}</h4>
                {this.renderUnreadChannels()}
            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
        fontFamily: 'Open Sans, sans-serif',
        overflowY: 'auto',
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
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(60px, 1fr))',
        gap: '10px',
        alignItems: 'stretch',
        justifyItems: 'stretch',
        padding: '0 0 1.8rem 0px',
    },
    contactItem: {
        display: 'flex',
        alignItems: 'center',
        padding: '8px 0',
        borderBottom: '1px solid #f0f0f0',
        transition: 'background-color 0.2s ease',
        cursor: 'pointer',
        ':hover': {backgroundColor: '#f0f8ff'},
    },
    contactCard: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: '100%',
        flexDirection: 'column',
        boxSizing: 'border-box',
    },
    viewTitle: {
        fontSize: '18px',
        margin: '0',
        textAlign: 'center',
        marginBottom: '.6rem',
        fontFamily: 'inherit',
        fontWeight: '600',
    },
    ballStatusStyle: {
        position: 'absolute',
        top: 'auto',
        right: '-4px',
        bottom: '-4px',
        width: '25px',
        height: '25px',
        color: 'white',
        display: 'inline-flex',
        borderRadius: '50%',
        zIndex: 5,
    },
    chatsTitle: {
        fontSize: '14px',
        fontFamily: 'inherit',
        fontWeight: '900',
    },
    contactCardTitle: {
        fontSize: '14px',
    },
    avatarCircle: {
        cursor: 'pointer',
        position: 'relative',
        width: '50px',
        height: '50px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    contactAvatar: {
        position: 'relative',
        marginRight: '12px',
    },
    avatar: {
        width: '100%',
        height: '100%',
        borderRadius: '50%',
        objectFit: 'cover',
    },
    avatarEmoji: {
        fontSize: '32px',
        display: 'block',
    },
    avatarEmoji2: {
        fontSize: '22px',
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
        padding: '8px',
        borderRadius: '3px',
        fontSize: '16px',
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
