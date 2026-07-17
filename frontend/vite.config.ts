// eslint-disable-next-line spaced-comment
/// <reference types="vitest/config" />
import fs from "fs";
import path from "path";
import {defineConfig} from "vite";
import {sveltekit} from "@sveltejs/kit/vite";
import {svelteTesting} from "@testing-library/svelte/vite";


const license = fs.readFileSync(path.join(__dirname, "src", "lib", "assets", "license.txt")).toString();

// https://vitejs.dev/config/
export default defineConfig({
    // svelteTesting adds test auto-cleanup and browser condition resolution,
    // both only active while running vitest.
    plugins: [sveltekit(), svelteTesting()],
    define: {
        __INSTALLER_LICENSE__: JSON.stringify(license)
    },
    test: {
        environment: "jsdom",
        setupFiles: ["src/setupTests.ts"],
        include: ["src/**/*.{test,spec}.{js,ts}", "src/**/*.{test,spec}.svelte"],
    }
});
