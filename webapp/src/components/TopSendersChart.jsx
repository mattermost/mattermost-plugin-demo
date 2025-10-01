import React from 'react';

const makeStyles = () => ({
    container: {
        width: '100%',
    },
    column: {
        display: 'flex',
        flexDirection: 'column',
        gap: 12,
    },
    row: {
        display: 'flex',
        alignItems: 'center',
        gap: 12,
    },
    name: {
        width: 90,
        fontSize: 13,
        color: '#3d3d3d',
        textAlign: 'right',
    },
    barBg: {
        flex: 1,
        background: '#e9eef5',
        borderRadius: 6,
        overflow: 'hidden',
    },
    bar: (width) => ({
        height: 18,
        width: `${width}%`,
        background: '#2f6feb',
        borderRadius: 6,
        transition: 'width 0.3s',
    }),
    value: {
        width: 48,
        fontSize: 12,
        color: '#555',
        textAlign: 'left',
    },
});

const data = [
    {name: 'Ana', conversations: 240},
    {name: 'Luis', conversations: 310},
    {name: 'María', conversations: 180},
    {name: 'Carlos', conversations: 355},
    {name: 'Sofía', conversations: 205},
];

const maxconversations = Math.max(...data.map((d) => d.conversations));

const TopSendersChart = () => {
    const styles = makeStyles();
    return (
        <div style={styles.container}>
            <div style={styles.column}>
                {data.map((d, i) => (
                    <div
                        key={i}
                        style={styles.row}
                    >
                        <div style={styles.name}>{d.name}</div>
                        <div style={styles.barBg}>
                            <div
                                style={styles.bar((d.conversations / maxconversations) * 100)}
                                title={`${d.conversations} mensajes`}
                            />
                        </div>
                        <div style={styles.value}>{d.conversations}</div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default TopSendersChart;
