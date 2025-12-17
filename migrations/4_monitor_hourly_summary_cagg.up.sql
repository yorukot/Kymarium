-- Ensure the pings table is a hypertable on the time column so continuous aggregates work.
SELECT create_hypertable('pings', 'time', if_not_exists => TRUE);

-- 30-minute rollup for availability and latency percentiles.
CREATE MATERIALIZED VIEW monitor_30min_summary
WITH (timescaledb.continuous) AS
SELECT
    monitor_id,
    region,
    time_bucket('30 minutes', time) AS bucket,
    count(*) AS total_count,
    count(*) FILTER (
        WHERE status = 'successful' AND latency <= 5000
    ) AS good_count,
    percentile_cont(0.50) WITHIN GROUP (ORDER BY latency) AS p50_ms,
    percentile_cont(0.75) WITHIN GROUP (ORDER BY latency) AS p75_ms,
    percentile_cont(0.90) WITHIN GROUP (ORDER BY latency) AS p90_ms,
    percentile_cont(0.95) WITHIN GROUP (ORDER BY latency) AS p95_ms,
    percentile_cont(0.99) WITHIN GROUP (ORDER BY latency) AS p99_ms
FROM pings
GROUP BY monitor_id, region, bucket
WITH NO DATA;

-- Keep the materialized view fresh:
-- refresh last 24 hours every 15 minutes,
-- skipping the most recent 30 minutes to avoid hot buckets.
SELECT add_continuous_aggregate_policy(
    'monitor_30min_summary',
    start_offset => INTERVAL '24 hours',
    end_offset   => INTERVAL '30 minutes',
    schedule_interval => INTERVAL '15 minutes'
);

-- Retention: drop raw ping data older than 90 days.
SELECT add_retention_policy('pings', INTERVAL '90 days');
