// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {ComponentType} from 'react';

export const CATEGORIES = [
    'actions',
    'forms',
    'feedback',
    'images',
    'layout',
    'progress',
    'status',
    'cards',
    'messaging',
    'navigation',
    'patterns',
    'banners',
    'chrome',
] as const;

export type Category = typeof CATEGORIES[number];

export const CATEGORY_LABELS: Record<Category, string> = {
    actions: 'Actions',
    forms: 'Forms and Input',
    feedback: 'Feedback and Notices',
    images: 'Images and Icons',
    layout: 'Layout and Containers',
    progress: 'Progress Indicators',
    status: 'Status Indicators',
    cards: 'Cards and Previews',
    messaging: 'Messaging',
    navigation: 'Navigation',
    patterns: 'Patterns',
    banners: 'Banners',
    chrome: 'Chrome',
};

export type PreviewSurface = 'sidebar';

export type CatalogEntry = {
    id: string;
    name: string;
    category: Category;
    description: string;
    preview: ComponentType;
    detail: ComponentType;
    previewSurface?: PreviewSurface;
};
