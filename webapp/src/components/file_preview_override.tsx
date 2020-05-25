import React from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

type Props = {
    fileInfo: FileInfo;
    onModalDismissed: () => void;
};

export default function FilePreviewOverride(props: Props) {
    return (
        <div>
            <h3>{props.fileInfo.name}</h3>
            <button onClick={props.onModalDismissed}>
                {'Close'}
            </button>
        </div>
    );
}
