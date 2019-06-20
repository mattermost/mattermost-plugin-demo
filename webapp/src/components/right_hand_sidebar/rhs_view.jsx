import React from 'react';

import {FormattedMessage} from 'react-intl';

export default class RHSView extends React.PureComponent {
    render() {
        return (
            <div style={style.rhs}>
                <br/>
                <br/>
                <br/>
                <br/>
                <FormattedMessage
                    id='rhs.triggered'
                    defaultMessage='You have triggered the right-hand sidebar component of the demo plugin.'
                />
                <br/>
                <br/>
                <br/>
                <br/>
                <FormattedMessage
                    id='demo.testintl'
                    defaultMessage='This is the default string'
                />
            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
    },
};