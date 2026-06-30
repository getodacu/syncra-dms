import type { AuthSession, AuthUser } from '$lib/server/auth';

declare global {
	namespace App {
		interface Locals {
			session: AuthSession | null;
			user: AuthUser | null;
			logger?: {
				info?: (message: string, attrs?: Record<string, unknown>) => void;
				error?: (message: string, attrs?: Record<string, unknown>) => void;
			};
		}
	}
}

export {};
