import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {NaiveUiResolver} from 'unplugin-vue-components/resolvers'

const proxyTarget: string = 'http://192.168.0.111:8080'

export default defineConfig({
    base: '/',
    server: {
        proxy: {
            '/api': {
                target: proxyTarget,
                changeOrigin: true,
            },
            '/ws': {
                target: proxyTarget,
                changeOrigin: true,
                ws: true
            }
        },
    },
    resolve: {
        extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    plugins: [
        vue(),
        AutoImport({
            imports: [
                'vue',
                {
                    'naive-ui': [
                        'useDialog',
                        'useMessage',
                        'useNotification',
                        'useLoadingBar'
                    ]
                }
            ]
        }),
        Components({
            resolvers: [NaiveUiResolver()]
        })
    ],
    build: {
        chunkSizeWarningLimit: 3000,
    }
})
