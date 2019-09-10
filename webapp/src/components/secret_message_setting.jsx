import React from 'react';
import PropTypes from 'prop-types';
const Setting = window.SystemConsoleSetting;

export default class SecretMessageSetting extends React.PureComponent {
    static propTypes = {
        id: PropTypes.string.isRequired,
        label: PropTypes.string.isRequired,
        helpText: PropTypes.node,
        value: PropTypes.any,
        disabled: PropTypes.bool.isRequired,
        config: PropTypes.object.isRequired,
        currentState: PropTypes.object.isRequired,
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
        }
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

    toggleSecretMessage = () => {
        this.setState({
            showSecretMessage: !this.state.showSecretMessage,
        });
    }

    handleChange = (e) => {
        this.props.onChange(this.props.id, e.target.value);
    }

    render() {
        return (
            <Setting
                label={this.props.label}
                inputId={this.props.id}
                helpText={this.props.helpText}
                setByEnv={this.props.setByEnv}
            >
                {this.state.showSecretMessage &&
                    <textarea
                        className="form-control"
                        rows={5}
                        value={this.props.value}
                        disabled={this.props.disabled || this.props.setByEnv}
                        onChange={this.handleChange}
                    />
                }
                <div className="help-text">
                    <button className="btn btn-default" onClick={this.toggleSecretMessage}>
                        {this.state.showSecretMessage && "Hide Secret Message"}
                        {!this.state.showSecretMessage && "Show Secret Message"}
                    </button>
                </div>
            </Setting>
        );
    }
}
