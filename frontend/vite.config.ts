import fs from "fs";
import path from "path";
import {defineConfig} from "vite";
import {sveltekit} from "@sveltejs/kit/vite";
import pkg from "./package.json";


const license = fs.readFileSync(path.join(__dirname, "src", "lib", "assets", "license.txt")).toString();

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [sveltekit()],
    define: {
        __INSTALLER_LICENSE__: JSON.stringify(license),
        __APP_VERSION__: JSON.stringify(pkg.version)
    }
});
