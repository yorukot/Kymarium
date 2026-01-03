import type { PublicTimelinePoint } from '$lib/types';

export const statusMeta = {
	up: { label: 'Operational', tone: 'text-success', dot: 'bg-success' },
	down: { label: 'Degraded', tone: 'text-destructive', dot: 'bg-destructive' }
};

export function formatSli(value?: number) {
	if (value === undefined || Number.isNaN(value)) return '—';
	return `${value.toFixed(2)}%`;
}

export function formatDate(value: string) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return value;
	return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

export function timelineTone(point: PublicTimelinePoint) {
	const total = point.success + point.fail;
	if (total === 0) return 'bg-muted';
	if (point.fail === 0) return 'bg-success';
	return 'bg-destructive';
}

export function formatUpTo2Decimals(n: number): string {
  // round to 3 decimals (mitigate float noise a bit)
  const r = Math.round((n + Number.EPSILON) * 1000) / 1000;

  // keep up to 3, then trim trailing zeros + optional dot
  return r.toFixed(2).replace(/\.?0+$/, "");
}