import '@mattermost/compass-ui/styles';

import Plugin from './plugin';
import {id as pluginId} from './manifest';

window.registerPlugin(pluginId, new Plugin());
