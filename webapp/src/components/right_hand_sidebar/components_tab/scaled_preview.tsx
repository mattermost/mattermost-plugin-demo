// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useLayoutEffect, useRef, useState} from 'react';

type Props = {
    children: React.ReactNode;
    className?: string;
};

export default function ScaledPreview({children, className}: Props) {
    const frameRef = useRef<HTMLDivElement>(null);
    const stageRef = useRef<HTMLDivElement>(null);
    const [scale, setScale] = useState(1);
    const [stageSize, setStageSize] = useState({width: 0, height: 0});

    useLayoutEffect(() => {
        const frame = frameRef.current;
        const stage = stageRef.current;
        if (!frame || !stage) {
            return undefined;
        }

        const fit = () => {
            const stageWidth = Math.max(stage.scrollWidth, stage.offsetWidth, 1);
            const stageHeight = Math.max(stage.scrollHeight, stage.offsetHeight, 1);
            const next = Math.min(
                frame.clientWidth / stageWidth,
                frame.clientHeight / stageHeight,
                1,
            );
            setStageSize({width: stageWidth, height: stageHeight});
            setScale(Number.isFinite(next) && next > 0 ? next : 1);
        };

        fit();

        if (typeof ResizeObserver === 'undefined') {
            return undefined;
        }

        const observer = new ResizeObserver(fit);
        observer.observe(frame);
        observer.observe(stage);
        return () => observer.disconnect();
    }, []);

    let previewClassName = 'DemoCompassShowcase__cardPreview';
    if (className) {
        previewClassName += ` ${className}`;
    }

    return (
        <div
            ref={frameRef}
            className={previewClassName}
        >
            <div
                className='DemoCompassShowcase__cardScaler'
                style={{
                    width: stageSize.width * scale,
                    height: stageSize.height * scale,
                }}
            >
                <div
                    ref={stageRef}
                    className='DemoCompassShowcase__cardStage'
                    style={{transform: `scale(${scale})`}}
                >
                    {children}
                </div>
            </div>
        </div>
    );
}
