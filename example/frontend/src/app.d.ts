// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			session: import("$lib/server/auth").AuthSession | null;
			user: import("$lib/server/auth").AuthUser | null;
			adminUser: import("$lib/server/auth").AuthUser | null;
			impersonation: import("$lib/server/auth").AuthImpersonation | null;
			logger: import("$lib/server/logging").Logger;
			requestId: string;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

declare module "*?worker" {
	const WorkerConstructor: {
		new (): Worker;
	};
	export default WorkerConstructor;
}

declare module "*?raw" {
	const content: string;
	export default content;
}

export {};
