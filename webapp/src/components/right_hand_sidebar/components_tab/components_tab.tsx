import React, {useEffect, useRef, useState} from 'react';

import Detail from './detail';
import Gallery from './gallery';
import type {CatalogEntry} from './types';

import './components_tab.css';

// Slightly longer than --duration-moderate so animationend normally wins.
const EXIT_FALLBACK_MS = 400;

export default function ComponentsTab() {
    const [selected, setSelected] = useState<CatalogEntry | null>(null);
    const [exiting, setExiting] = useState(false);
    const triggerRef = useRef<HTMLButtonElement | null>(null);

    const finishExit = () => {
        setSelected(null);
        setExiting(false);
    };

    const handleBack = () => {
        if (window.matchMedia?.('(prefers-reduced-motion: reduce)').matches) {
            finishExit();
            return;
        }
        setExiting(true);
    };

    // Safety net so the detail pane can never get stuck if animationend doesn't fire.
    useEffect(() => {
        if (!exiting) {
            return undefined;
        }
        const timeout = window.setTimeout(finishExit, EXIT_FALLBACK_MS);
        return () => window.clearTimeout(timeout);
    }, [exiting]);

    // Restore focus once, when the detail pane is gone, unless the user has already moved it.
    useEffect(() => {
        if (selected) {
            return;
        }
        const active = document.activeElement;
        if (triggerRef.current && (!active || active === document.body)) {
            triggerRef.current.focus({preventScroll: true});
        }
    }, [selected]);

    return (
        <div className='DemoCompassShowcase__nav'>
            <div
                className={`DemoCompassShowcase__pane DemoCompassShowcase__pane--base${selected && !exiting ? ' DemoCompassShowcase__pane--pushed' : ''}`}
                inert={selected !== null}
                aria-hidden={selected ? true : undefined}
            >
                <Gallery
                    onSelect={(entry, trigger) => {
                        triggerRef.current = trigger;
                        setSelected(entry);
                    }}
                />
            </div>
            {selected && (
                <div
                    className={`DemoCompassShowcase__pane DemoCompassShowcase__pane--push${exiting ? ' DemoCompassShowcase__pane--exiting' : ''}`}
                    inert={exiting}
                    onAnimationEnd={(event) => {
                        if (exiting && event.target === event.currentTarget) {
                            finishExit();
                        }
                    }}
                >
                    <Detail
                        key={selected.id}
                        entry={selected}
                        onBack={handleBack}
                    />
                </div>
            )}
        </div>
    );
}
