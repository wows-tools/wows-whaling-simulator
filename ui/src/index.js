import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import Error from "./components/Error";
import reportWebVitals from "./reportWebVitals";
import { createBrowserRouter, RouterProvider, Route } from "react-router-dom";
import LootboxList from "./components/LootboxList";
import Lootbox from "./components/Lootbox";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <Error />,
    children: [
      {
        path: "",
        element: <LootboxList />,
      },
      {
        path: "lootboxes/:lootboxId",
        element: <Lootbox />,
      },
    ],
  },
]);

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<RouterProvider router={router} />);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
