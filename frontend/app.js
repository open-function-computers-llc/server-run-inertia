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

// Listen for messages on the hot reload websocket
const inertiaPageData = JSON.parse(document.getElementById("app").getAttribute("data-page"));
const liveReloadPort = inertiaPageData.props.liveReloadPort;

if (liveReloadPort && liveReloadPort != "") {
    const socket = new WebSocket(`ws://${window.location.hostname}:${liveReloadPort}/frontend-reload`);
    socket.onmessage = (event) => {
        console.log("[Live Reload]", event.data);
        if (event.data === "reload") {
            window.location.reload();
        }
    }
}

