// Note: The fs.Stats deprecation warning is from Vue CLI's webpack dependencies
// which are not yet compatible with Node.js v23. This is suppressed in the build
// script. All user dependencies have been updated to their latest versions.
// Consider migrating to Vite for better Node.js v23+ compatibility.
const webpack = require('webpack')

module.exports = {
	filenameHashing: false,
	productionSourceMap: false,
	outputDir: '../static',
	publicPath: '/',
	configureWebpack: {
		plugins: [
			new webpack.DefinePlugin({
				__VUE_I18N_FULL_INSTALL__: JSON.stringify(true),
				__VUE_I18N_LEGACY_API__: JSON.stringify(false),
				__INTLIFY_PROD_DEVTOOLS__: JSON.stringify(false),
				__INTLIFY_DROP_MESSAGE_COMPILER__: JSON.stringify(false),
			}),
		],
	},
	devServer: {
		port: 8081,
		https: false,
		client: {
			webSocketURL:'auto://0.0.0.0/ws'
		},
		proxy: {
			'^/api|^/css|^/oicd': {
				target: "http://localhost:8080",
				changeOrigin: true,
				secure: false,
			}
		}
	}
}
