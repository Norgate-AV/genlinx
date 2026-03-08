import { defineConfig } from "tsup";

export default defineConfig({
    entry: ["src/app.ts"],
    dts: false,
    format: ["esm"],
    target: "node22",
    minify: "terser",
    minifyWhitespace: true,
    minifyIdentifiers: true,
    minifySyntax: true,
    treeshake: true,
    sourcemap: false,
    splitting: true,
});
