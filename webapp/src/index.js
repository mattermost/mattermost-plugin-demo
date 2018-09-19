import Plugin from './plugin';
import Manifest from './manifest';

window.registerPlugin(Manifest.PluginId, new Plugin());
