import React from 'react';

import {FileInfo} from 'mattermost-redux/types/files';
import {Theme} from 'mattermost-redux/types/preferences';

type Props = {
    fileInfo: FileInfo;
    onModalDismissed: () => void;
    theme: Theme;
};

export default function FilePreviewOverride(props: Props) {
    const {theme} = props;

    const style = {
        backgroundColor: theme.centerChannelBg,
        color: theme.centerChannelColor,
        paddingTop: '10px',
        paddingBottom: '10px',
    };

    const buttonStyle = {
        backgroundColor: theme.buttonBg,
        color: theme.buttonColor,
    };

    return (
        <div style={style}>
            <h3>{props.fileInfo.name}</h3>
            <button
                className={'save-button btn'}
                style={buttonStyle}
                onClick={props.onModalDismissed}
            >
                {'Close'}
            </button>
        </div>
    );
}
