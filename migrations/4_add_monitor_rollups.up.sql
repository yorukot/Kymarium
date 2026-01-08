CREATE MATERIALIZED VIEW monitor_2min_summary
WITH (timescaledb.continuous) AS
SELECT
    monitor_id,
    region_id,
    time_bucket('2 minutes', time) AS bucket,
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
GROUP BY monitor_id, region_id, bucket
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
    'monitor_2min_summary',
    start_offset => INTERVAL '24 hours',
    end_offset   => INTERVAL '2 minutes',
    schedule_interval => INTERVAL '1 minute'
);

SELECT add_retention_policy('monitor_2min_summary', INTERVAL '24 hours');

ALTER MATERIALIZED VIEW monitor_2min_summary
SET (timescaledb.materialized_only = false);

CREATE MATERIALIZED VIEW monitor_10min_summary
WITH (timescaledb.continuous) AS
SELECT
    monitor_id,
    region_id,
    time_bucket('10 minutes', time) AS bucket,
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
GROUP BY monitor_id, region_id, bucket
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
    'monitor_10min_summary',
    start_offset => INTERVAL '7 days',
    end_offset   => INTERVAL '10 minutes',
    schedule_interval => INTERVAL '5 minutes'
);

SELECT add_retention_policy('monitor_10min_summary', INTERVAL '7 days');

ALTER MATERIALIZED VIEW monitor_10min_summary
SET (timescaledb.materialized_only = false);
