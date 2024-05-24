import {defineConfig} from 'vite'
import {resolve} from 'path'
import react from '@vitejs/plugin-react'
// import basicSsl from '@vitejs/plugin-basic-ssl'

/**
 * @description
 * @param command { "build" | "server" } 模式
 * @param mode { string } 模式
 * @param isSsrBuild { string } 是否为SSR 构建
 * @param isPreview { string } 是否为正在预览的构建产物
 * @returns
 */
export default defineConfig(async ({command}) => {
	if (command === 'serve') {
		return {
			// dev 独有配置
			plugins: [
				react(),
				// 				basicSsl({
				// 					/** name of certification */
				// 					name: 'test',
				// 					/** custom trust domains */
				// 					domains: ['*.custom.com'],
				// 					/** custom certification directory */
				// 					certDir: '/Users/.../.devServer/cert',
				// 				}),
			],

			resolve: {
				alias: {
					'@': resolve(__dirname, 'src'),
				},
			},

			// Vite options tailored for Tauri development and only applied in `tauri dev` or `tauri build`
			//
			// 1. prevent vite from obscuring rust errors
			clearScreen: false,
			// 2. tauri expects a fixed port, fail if that port is not available
			server: {
				host: '0.0.0.0',
				port: 3000,
				strictPort: true,
				watch: {
					// 3. tell vite to ignore watching `src-tauri`
					ignored: ['**/src-tauri/**'],
				},
			},
		}
	} else {
		return {
			// build 独有配置
			// 开发或生产环境服务的公共基础路径。合法的值包括以下几种：
			// 绝对 URL 路径名，例如 /foo/
			// 完整的 URL，例如 https://foo.com/（原始的部分在开发环境中不会被使用）
			// 空字符串或 ./（用于嵌入形式的开发）
			base: './',
			mode: 'production',
			server: {
				host: '0.0.0.0',
				port: 443,
				strictPort: true,
			},
			assetsDir: 'public',
			chunkSizeWarningLimit: 1000,
			// 配置打包文件路径和命名
			minify: 'esbuild',
			outDir: 'dist',
			// 取消计算文件大小，加快打包速度
			reportCompressedSize: false,
			sourcemap: false,
			target: 'esnext',
			terserOptions: {
				compress: {
					// 生产环境时移除console.log调试代码
					drop_console: true,
					drop_debugger: true,
				},
			},
			rollupOptions: {
				// // 打包时忽略某些包，避免打包时间过长
				// 				external: ['react'],
			},
		}
	}
})
