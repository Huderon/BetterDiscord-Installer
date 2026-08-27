import {node} from "@zerebos/eslint-config";
import ts from "@zerebos/eslint-config-typescript";
import {defineConfig} from "eslint/config";
import {build} from "@zerebos/eslint-config-svelte";
import tseslint from "typescript-eslint";
import svelteConfig from "./svelte.config.js";


/** @type {import("@zerebos/eslint-config-typescript").ConfigArray} */
export default defineConfig(
    defineConfig(
        ...node,
        ...ts.configs.recommendedWithTypes
    ),
    ...build(svelteConfig),
    {
        languageOptions: {
            globals: {
                __INSTALLER_LICENSE__: "readonly",
            }
        }
    },
    {
        // Root-level JS config files aren't part of the tsconfig project, and
        // relying on the project service's default-project fallback for them
        // is flaky in CI. They don't benefit from type-aware rules anyway.
        "files": ["*.js"],
        "extends": [tseslint.configs.disableTypeChecked],
    },
    {
        ignores: ["build/", ".svelte-kit/", "dist/", "src/lib/wailsjs/", "node_modules/"]
    },
);
