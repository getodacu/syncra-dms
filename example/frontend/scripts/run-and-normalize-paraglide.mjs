import { spawn } from "node:child_process";
import { constants as osConstants } from "node:os";
import { normalizeParaglide } from "./normalize-paraglide.mjs";

const [command, ...args] = process.argv.slice(2);

if (!command) {
	console.error("Usage: node scripts/run-and-normalize-paraglide.mjs <command> [...args]");
	process.exit(1);
}

let exiting = false;
let childExited = false;
let forceKillTimer;
let normalized = false;

const child = spawn(command, args, {
	stdio: "inherit",
	shell: process.platform === "win32"
});

child.on("error", (error) => {
	console.error(error.message);
	normalizeAndExit(1);
});

child.on("exit", (status, signal) => {
	childExited = true;
	normalizeAndExit(status, signal);
});

for (const signal of ["SIGINT", "SIGTERM", "SIGHUP"]) {
	process.once(signal, () => {
		normalizeOnce();

		if (!childExited) {
			child.kill(signal);
		}

		forceKillTimer = setTimeout(() => {
			if (!childExited) {
				child.kill("SIGKILL");
			}
		}, 5000);
		forceKillTimer.unref();
	});
}

process.on("exit", () => {
	normalizeOnce();
});

function normalizeOnce() {
	if (normalized) return;
	normalized = true;
	normalizeParaglide();
}

function normalizeAndExit(status, signal) {
	if (exiting) return;
	exiting = true;

	if (forceKillTimer) {
		clearTimeout(forceKillTimer);
	}

	try {
		normalizeOnce();
	} catch (error) {
		console.error(error instanceof Error ? error.message : error);
		process.exit(1);
	}

	if (signal) {
		process.exit(128 + (osConstants.signals[signal] ?? 1));
	}

	process.exit(status ?? 1);
}
