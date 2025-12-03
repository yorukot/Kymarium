-- Add ping monitor support by extending the monitor_type enum.
ALTER TYPE "monitor_type" ADD VALUE IF NOT EXISTS 'ping';
