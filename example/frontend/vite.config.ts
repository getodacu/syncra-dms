import { paraglideVitePlugin } from '@inlang/paraglide-js';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { normalizeParaglide } from './scripts/normalize-paraglide.mjs';

function normalizeParaglideOutputPlugin() {
	return {
		name: 'normalize-paraglide-output',
		buildStart() {
			normalizeParaglide();
		},
		watchChange() {
			normalizeParaglide();
		}
	};
}

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide',
			strategy: ['cookie', 'baseLocale'],
			emitTsDeclarations: true,
			emitGitIgnore: false,
			outputStructure: 'locale-modules'
		}),
		normalizeParaglideOutputPlugin()
	]
});
