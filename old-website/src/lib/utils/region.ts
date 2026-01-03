import type { Region } from '$lib/types';

export function regionFlagIcon(region: Region): string | null {
	const countryCode = region.name.split('-')[0]?.trim().toLowerCase();
	return countryCode ? `cif:${countryCode}` : null;
}

export function regionFlagEmoji(region: Region): string | null {
	// 假設 name 是 "us-east-1" / "sg-sin-1"
	const countryCode = region.name.split('-')[0];
	return countryCodeToFlagEmoji(countryCode);
}

export function countryCodeToFlagEmoji(code?: string): string | null {
	if (!code) return null;

	const cc = code.trim().toUpperCase();
	if (cc.length !== 2) return null;

	const A = 0x1f1e6; // 🇦
	const base = 'A'.charCodeAt(0);

	return String.fromCodePoint(A + cc.charCodeAt(0) - base, A + cc.charCodeAt(1) - base);
}
