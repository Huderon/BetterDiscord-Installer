import adapter from "@sveltejs/adapter-static";
import {vitePreprocess} from "@sveltejs/vite-plugin-svelte";


/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: vitePreprocess(),
    kit: {
        adapter: adapter({fallback: "index.html"}),
        alias: {
            "@assets/*": "./src/lib/assets/*",
            "@wails/*": "./src/lib/wailsjs/runtime/*",
            "@api": "./src/lib/wailsjs/go/api/Controller",
            "@app": "./src/lib/wailsjs/go/main/App",
            "@models": "./src/lib/wailsjs/go/models.ts"
        }
    }
};

export default config;
