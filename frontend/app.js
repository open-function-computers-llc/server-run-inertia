import { createApp, h } from "vue";
import { createInertiaApp, Link } from "@inertiajs/vue3";
import * as bootstrap from 'bootstrap';

createInertiaApp({
    title: (title) => title ? `${title} | Server Administration` : "Server Administration",
    resolve: (name) => require(`./Pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
        .use(plugin)
        .mount(el);
    }
})
