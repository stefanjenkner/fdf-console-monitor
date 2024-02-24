/** @type {import('ts-jest').JestConfigWithTsJest} */
module.exports = {
  roots: [
    "<rootDir>/src"
  ],
  transform: {
    '^.+\\.(ts|tsx)?$': 'ts-jest',
  },
  testEnvironment: 'node',
};