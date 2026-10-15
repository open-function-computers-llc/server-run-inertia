let path = require("path");
let mix = require('laravel-mix');

mix.alias({
   "@": path.join(__dirname, "frontend"),
});

mix.js('frontend/app.js', 'dist')
   .version()
   .sourceMaps()
   .vue()
   .setPublicPath('dist');

mix.sass('frontend/scss/app.scss', 'dist')
   .version()
   .setPublicPath('dist');

mix.copy("server/index.html", "dist/index.html");

// node-notifier's bundled terminal-notifier is x86-only and fails on Apple Silicon (spawn error -86)
mix.disableNotifications();
