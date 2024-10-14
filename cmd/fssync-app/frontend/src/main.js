import { createApp } from 'vue'
import App from './App.vue'
import './style.css';

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import FileIcon from "@/components/FileIcon/index.vue"
import router from './router'
import download from './plugins/download' // plugins

const app = createApp(App);

app.use(router)
    .use(ElementPlus)
    .component('file-icon', FileIcon)
    .mount('#app')

app.config.globalProperties.$download = download