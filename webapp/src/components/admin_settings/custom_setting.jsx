import React from 'react';
import PropTypes from 'prop-types';

export default class CustomSetting extends React.PureComponent {
    static propTypes = {
        id: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        helpText: PropTypes.node,
        value: PropTypes.any,
        disabled: PropTypes.bool.isRequired,
        config: PropTypes.object.isRequired,
        license: PropTypes.object.isRequired,
        setByEnv: PropTypes.bool.isRequired,
        onChange: PropTypes.func.isRequired,
        registerSaveAction: PropTypes.func.isRequired,
        setSaveNeeded: PropTypes.func.isRequired,
        unRegisterSaveAction: PropTypes.func.isRequired,
    }

    constructor(props) {
        super(props);

        this.state = {
            value: this.props.value || 0,
        };
    }

    handleChange = (e) => {
        this.setState({value: e.target.value});

        this.props.onChange(this.props.id, e.target.value);
        this.props.setSaveNeeded();
    }

    render() {
        return (
            <div style={style.div}>
                <div style={style.text}>
                    {'Demo Custom Setting'}</div>
                <div style={style.text}>
                    {this.state.value}</div>
                <input
                    type='range'
                    min='0'
                    max='10'
                    disabled={this.props.disabled}
                    value={this.state.value}
                    onChange={this.handleChange}
                />
            </div>
        );
    }
}

const style = {
    div: {
        padding: '10px 50px 15px 50px',
        border: '1px solid #ddd',
        margin: '10px',
    },
    text: {
        padding: '5px',
        textAlign: 'center',
        fontWeight: '900',
    },
    input: {
        margin: '5px',
    },
};
