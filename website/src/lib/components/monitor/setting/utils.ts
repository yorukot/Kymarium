// Utility helpers for monitor settings payload shaping.

/**
 * Convert a multi-line header string into a header record.
 * Each line should look like "Key: Value".
 */
export function parseHeaders(raw: string | undefined): Record<string, string> | undefined {
	if (!raw) return undefined;

	const lines = raw
		.split('\n')
		.map((line) => line.trim())
		.filter(Boolean);

	if (!lines.length) return undefined;

	const headers: Record<string, string> = {};
	for (const line of lines) {
		const idx = line.indexOf(':');
		if (idx === -1) continue;
		const key = line.slice(0, idx).trim();
		const value = line.slice(idx + 1).trim();
		if (key && value) headers[key] = value;
	}

	return Object.keys(headers).length ? headers : undefined;
}

/**
 * Normalize accepted status code selections into unique numeric codes.
 * Numeric entries (e.g. "200", "404") are kept.
 * Shorthand ranges like "2xx" are expanded to the real, registered codes for that class.
 */
export function normalizeStatusCodes(codes: string[]): number[] {
	const result = new Set<number>();

	for (const code of codes) {
		const range = /^([1-5])xx$/.exec(code);
		if (range) {
			const cls = Number(range[1]);
			ALL_STATUS_CODES.forEach((value) => {
				if (Math.floor(value / 100) === cls) result.add(value);
			});
			continue;
		}

		const parsed = Number(code);
		if (Number.isInteger(parsed) && ALL_STATUS_CODES.has(parsed)) {
			result.add(parsed);
		}
	}

	return Array.from(result).sort((a, b) => a - b);
}

/**
 * Replace full numeric coverage of a class with its shorthand token (e.g. 2xx),
 * and drop numeric codes when the shorthand is already selected.
 */
export function condenseStatusTokens(tokens: string[]): string[] {
	const next = new Set(tokens);

	for (let cls = 1; cls <= 5; cls++) {
		const classToken = `${cls}xx`;
		const classCodes = CLASS_STATUS_CODES[cls] ?? [];
		const classCodeStrings = classCodes.map(String);

		// If shorthand already selected, remove all numeric members.
		if (next.has(classToken)) {
			classCodeStrings.forEach((code) => next.delete(code));
			continue;
		}

		// If all numeric codes are present, collapse to shorthand.
		const hasAll = classCodeStrings.length > 0 && classCodeStrings.every((code) => next.has(code));
		if (hasAll) {
			classCodeStrings.forEach((code) => next.delete(code));
			next.add(classToken);
		}
	}

	// Preserve deterministic order: numeric ascending, then class tokens ascending.
	return Array.from(next).sort((a, b) => {
		const numA = Number(a);
		const numB = Number(b);
		const isNumA = Number.isInteger(numA);
		const isNumB = Number.isInteger(numB);
		if (isNumA && isNumB) return numA - numB;
		if (isNumA) return -1;
		if (isNumB) return 1;
		return a.localeCompare(b);
	});
}

/**
 * Collapse a full set of status codes for a class back into its shorthand token (e.g. 2xx).
 * This keeps forms compact when all codes of a class are present.
 */
export function collapseStatusCodesToRanges(codes: number[]): string[] {
	const remaining = new Set(codes);
	const result: string[] = [];

	for (let cls = 1; cls <= 5; cls++) {
		const classCodes = CLASS_STATUS_CODES[cls] ?? [];
		const hasAll = classCodes.every((code) => remaining.has(code));
		if (hasAll && classCodes.length > 0) {
			result.push(`${cls}xx`);
			classCodes.forEach((code) => remaining.delete(code));
		}
	}

	// add any leftover specific codes
	result.push(
		...Array.from(remaining)
			.sort((a, b) => a - b)
			.map(String)
	);

	return result;
}

// --- Internal helpers ----

// Source: IANA/MDN registered status codes (HTTP/1.1 + common extensions).
export const ALL_STATUS_CODES = new Set<number>([
	// 1xx Informational
	100, 101, 102, 103,
	// 2xx Success
	200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
	// 3xx Redirection
	300, 301, 302, 303, 304, 305, 306, 307, 308,
	// 4xx Client Error
	400, 401, 402, 403, 404, 405, 406, 407, 408, 409,
	410, 411, 412, 413, 414, 415, 416, 417, 418, 421, 422,
	423, 424, 425, 426, 428, 429, 431, 451,
	// 5xx Server Error
	500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511
]);

const CLASS_STATUS_CODES: Record<number, number[]> = {};
for (const code of ALL_STATUS_CODES) {
	const cls = Math.floor(code / 100);
	if (!CLASS_STATUS_CODES[cls]) CLASS_STATUS_CODES[cls] = [];
	CLASS_STATUS_CODES[cls].push(code);
}
