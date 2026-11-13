# AGENTS.md

Demo plugin: exercises as much of the Mattermost plugin API as possible, so it's a reference other plugins copy from. The Go server is in `server/`, the webapp in `webapp/`.

## Commands

```bash
go test ./server/...                               # server tests
cd webapp && npx jest src/components/sample_flow   # webapp tests for one folder
cd webapp && npm run lint && npm run check-types
make check-style                                   # Go and webapp lint
make test                                          # full check before pushing
make deploy                                        # build and install on a local server
```

## Rules

- `min_server_version` is 12.0.0. Don't use plugin API methods or host webapp exports newer than that without bumping it.
- React, ReactDOM, `react/jsx-runtime`, react-intl, Redux, react-redux, react-router-dom and prop-types are webpack `externals` served by the host. Never bundle them or move them to `dependencies`; their devDependency versions match the host for type checking and tests.
- Every new registration or slash command gets an entry in `webapp/src/components/right_hand_sidebar/surfaces.js`, which the RHS Guide tab lists.
- New webapp code is TypeScript. Build UI from `@mattermost/compass-ui` subpath imports (`@mattermost/compass-ui/components/button`) and icons from `@mattermost/compass-icons/components/*`.
- Compass `Modal` and `PopoverMenu` don't supply a portal, backdrop or Escape handling. Wrap modals in `ModalLayer` (`webapp/src/components/modal_layer`) and pass `closeOnEscape={false}`; `specimens/live_overlays.tsx` shows the menu pattern.
- CSS class names start with `Demo` or `demo-`, because style-loader injects them globally next to the host's and other plugins' styles.
- User-facing strings use react-intl with an `id` and `defaultMessage`. Add new ids to `webapp/i18n/en.json` only; other locales come from translators.
- Write webapp tests with React Testing Library and `tests/render.tsx`, querying by role and accessible name.
