const { defineConfig } = require('@vue/cli-service');
const WebpackObfuscator = require('webpack-obfuscator');

module.exports = defineConfig({
  transpileDependencies: ['vuetify'],
  productionSourceMap: process.env.NODE_ENV === 'production' ? false : true,
  configureWebpack: {
    plugins:
      process.env.NODE_ENV === 'production'
        ? [
            new WebpackObfuscator(
              {
                compact: true,
                disableConsoleOutput: true,
                identifierNamesGenerator: 'mangled-shuffled',
                debugProtection: true,
                debugProtectionInterval: 4000
              },
              ['**chunk**']
            )
          ]
        : []
  }
});
