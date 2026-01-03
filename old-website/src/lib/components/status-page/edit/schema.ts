import type { StatusPageElement } from '$lib/types';
import { z } from 'zod';

export const Int64String = z.string().regex(/^-?\d+$/, 'Expected int64 string');

export const statusPageElementTypeSchema = z.enum([
	'historical_timeline',
	'current_status_indicator'
]);

export const statusPageMonitorInputSchema = z.object({
	id: Int64String.optional(),
	monitorId: Int64String,
	groupId: Int64String.optional(),
	name: z.string().min(1).max(255),
	type: statusPageElementTypeSchema,
	sortOrder: z.coerce.number().int().min(1)
});

export const statusPageElementInputSchema = z
	.object({
		id: Int64String.optional(),
		name: z.string().min(1).max(255),
		type: statusPageElementTypeSchema,
		sortOrder: z.coerce.number().int().min(1),
		monitor: z.coerce.boolean().default(false),
		monitorId: Int64String.optional(),
		monitors: z.array(statusPageMonitorInputSchema).default([])
	})
	.superRefine((value, ctx) => {
		if (value.monitor && !value.monitorId) {
			ctx.addIssue({
				code: 'custom',
				message: 'monitorId is required for monitor elements',
				path: ['monitorId']
			});
		}
	});

export const statusPageUpsertRequestSchema = z.object({
	name: z.string().min(1).max(255),
	slug: z.string().min(3).max(255),

	elements: z.array(statusPageElementInputSchema).default([])
});

export type StatusPageUpsertRequest = z.infer<typeof statusPageUpsertRequestSchema>;
export type StatusPageElementInput = z.infer<typeof statusPageElementInputSchema>;
export type StatusPageMonitorInput = z.infer<typeof statusPageMonitorInputSchema>;
export type StatusPageElementType = z.infer<typeof statusPageElementTypeSchema>;

export type StatusPageUpsertValues = {
	name: string;
	slug: string;
	elements: StatusPageElement[];
};
