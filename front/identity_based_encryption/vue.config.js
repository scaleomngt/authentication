module.exports = {
    devServer: {
        host: '127.0.0.1',
        port: 9992,
        proxy: {
            '/api': {
                target: `http://172.18.10.46:7777`,
                pathRewrite: { '^/api': '' },
            }
        }
    },
    productionSourceMap: false,
}