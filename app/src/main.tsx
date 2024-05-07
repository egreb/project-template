import { render } from "preact";
import "./index.css";
import { Routes } from "@generouted/react-router";

render(<Routes />, document.getElementById("app") as HTMLElement);
