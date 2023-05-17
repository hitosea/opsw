import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {NaiveUiResolver} from 'unplugin-vue-components/resolvers'
import autoprefixer from 'autoprefixer'

export default defineConfig({
    base: '/',
    server: {
        proxy: {
            '/api': {
                target: 'http://192.168.100.22:8080',
                changeOrigin: true,
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
    },
    css: {
        postcss: {
            plugins: [
                autoprefixer
            ]
        }
    }
})
