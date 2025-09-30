import React from 'react';

import WeeklyMessagesChart from './WeeklyMessagesChart';
import TopSendersChart from './TopSendersChart';

const ChartsDialog = () => {
    const styles = makeStyles();
    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <div style={styles.title}>{'Estadísticas de mensajes'}</div>
            </div>
            <div style={styles.grid}>
                <div style={styles.card}>
                    <div style={styles.cardTitle}>{'Cantidad de mensajes por empleado'}</div>
                    <TopSendersChart/>
                </div>
                <div style={styles.card}>
                    <div style={styles.cardTitle}>{'Mensajes de los últimos 7 días'}</div>
                    <WeeklyMessagesChart/>
                </div>
            </div>
        </div>
    );
};
const makeStyles = () => ({
    container: {
        display: 'flex',
        flexDirection: 'column',
        gap: 16,
        width: '100%',
        maxWidth: 840,
    },
    header: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
    },
    title: {
        fontSize: 18,
        fontWeight: 600,
    },
    grid: {
        display: 'grid',
        gridTemplateColumns: '1fr',
        gap: 16,
    },
    card: {
        border: '1px solid #e5e7eb',
        borderRadius: 8,
        padding: 16,
        background: '#fff',
    },
    cardTitle: {
        fontSize: 14,
        fontWeight: 600,
        marginBottom: 8,
    },
});

export default ChartsDialog;
