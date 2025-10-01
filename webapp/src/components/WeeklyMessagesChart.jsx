import React from 'react';

const data = [
    {week: 'Lun', messages: 120},
    {week: 'Mar', messages: 98},
    {week: 'Mié', messages: 150},
    {week: 'Jue', messages: 87},
    {week: 'Vie', messages: 175},
    {week: 'Sáb', messages: 60},
    {week: 'Dom', messages: 92},
];

const maxMessages = Math.max(...data.map((d) => d.messages));

const makeStyles = () => ({
    container: {
        width: '100%',
        padding: 12,
    },
    chart: {
        display: 'flex',
        alignItems: 'flex-end',
        height: 180,
        gap: 14,
    },
    barWrapper: {
        textAlign: 'center',
        flex: 1,
    },
    bar: (height) => ({
        height: `${height}px`,
        background: '#14a27f',
        borderRadius: 6,
        marginBottom: 8,
        transition: 'height 0.3s',
        boxShadow: '0 1px 2px rgba(0,0,0,0.08) inset',
    }),
    day: {
        fontSize: 12,
        color: '#3d3d3d',
    },
    value: {
        fontSize: 12,
        color: '#6b7280',
    },
});

const WeeklyMessagesChart = () => {
    const styles = makeStyles();
    return (
        <div style={styles.container}>
            <div style={styles.chart}>
                {data.map((d, i) => (
                    <div key={i} style={styles.barWrapper}>
                        <div
                            style={styles.bar((d.messages / maxMessages) * 140)}
                            title={`${d.messages} mensajes`}
                        />
                        <div style={styles.day}>{d.week}</div>
                        <div style={styles.value}>{d.messages}</div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default WeeklyMessagesChart;
