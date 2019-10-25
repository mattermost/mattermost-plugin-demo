import React from 'react';
import PropTypes from 'prop-types';

export default class SecretMessageSetting extends React.PureComponent {
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
            showSecretMessage: false,
        };
    }

    componentDidMount() {
        this.props.registerSaveAction(this.handleSave);
    }

    componentWillUnmount() {
        this.props.unRegisterSaveAction(this.handleSave);
    }

    handleSave = async () => {
        this.setState({
            error: '',
        });

        let error;
        return {error};
    }

    showSecretMessage = () => {
        this.setState({
            showSecretMessage: true,
        });
    }

    toggleSecretMessage = (e) => {
        e.preventDefault();

        this.setState({
            showSecretMessage: !this.state.showSecretMessage,
        });
    }

    handleChange = (e) => {
        this.props.onChange(this.props.id, e.target.value);
    }

    render() {
        return (
            <React.Fragment>
                {this.state.showSecretMessage &&
                    <textarea
                        style={style.input}
                        className='form-control input'
                        rows={5}
                        value={this.props.value}
                        disabled={this.props.disabled || this.props.setByEnv}
                        onInput={this.handleChange}
                    />
                }
                <div>
                    <button
                        className='btn btn-default'
                        onClick={this.toggleSecretMessage}
                        disabled={this.props.disabled}
                    >
                        {this.state.showSecretMessage && 'Hide Secret Message'}
                        {!this.state.showSecretMessage && 'Show Secret Message'}
                    </button>
                </div>
            </React.Fragment>
        );
    }
}

const style = {
    input: {
        marginBottom: '5px',
    },
};
