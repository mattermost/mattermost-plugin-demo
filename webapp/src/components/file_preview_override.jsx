import React from 'react';

export default function FilePreviewOverride(props) {
    return (
        <div>
            <h3>{props.fileInfo.name}</h3>
            <button onClick={props.onModalDismissed}>
                {'Close'}
            </button>
        </div>
    )
}
