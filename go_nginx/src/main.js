import "normalize.css";
import "skeleton-css/css/skeleton.css";
import "./global.scss";
import HMR from "@roxi/routify/hmr";
import App from "./App.svelte";

const app = HMR(App, { target: document.body }, "routify-app");

export default app;
