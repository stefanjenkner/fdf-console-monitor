const tsPlugin = require("@typescript-eslint/eslint-plugin");
const tsParser = require("@typescript-eslint/parser");
const stylisticTs = require("@stylistic/eslint-plugin-ts");

const tsOverrideConfig = tsPlugin.configs["eslint-recommended"].overrides[0];
const tsRecommemdedConfig = tsPlugin.configs.recommended;
const files = ["**/*.ts", "**/*.tsx"];

module.exports = [
  "eslint:recommended",
  {
    files,
    linterOptions: {
      reportUnusedDisableDirectives: true,
    },
    languageOptions: {
      parser: tsParser,
    },
    plugins: {
      "@typescript-eslint": tsPlugin,
      "@stylistic/ts": stylisticTs
    },
    rules: {
        "sort-imports": ["error"],
        "@stylistic/ts/indent": ["error", 4],
        "@stylistic/ts/quotes": ["error", "single"],
        "@stylistic/ts/semi": ["error", "always"],
    }
  },
  { files, rules: tsOverrideConfig.rules },
  { files, rules: tsRecommemdedConfig.rules },
];