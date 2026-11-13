import React, {type ReactNode} from 'react';

import './plugin_page.css';

type Props = {
    title: ReactNode;
    description?: ReactNode;
    nav?: ReactNode;
    children?: ReactNode;
};

export default function PluginPage({title, description, nav, children}: Props) {
    return (
        <main className='demo-page'>
            <header className='demo-page__header'>
                <h1 className='demo-page__title'>{title}</h1>
                {description && <p className='demo-page__description'>{description}</p>}
            </header>
            {nav && <div className='demo-page__nav'>{nav}</div>}
            <div className='demo-page__body'>{children}</div>
        </main>
    );
}

type FactsProps = {
    items: Array<{label: string; value: ReactNode}>;
};

export function Facts({items}: FactsProps) {
    return (
        <dl className='demo-page__facts'>
            {items.map(({label, value}) => (
                <div
                    key={label}
                    className='demo-page__fact'
                >
                    <dt>{label}</dt>
                    <dd><code>{value}</code></dd>
                </div>
            ))}
        </dl>
    );
}
