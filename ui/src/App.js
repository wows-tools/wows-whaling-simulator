import {
  Flex,
  View,
  Provider,
  defaultTheme,
  Switch,
  Text,
} from "@adobe/react-spectrum";
import React, { useState } from "react";
import AppFooter from "./components/Footer";
import AppHeader from "./components/Header";
import LootboxList from "./components/LootboxList";
import { Outlet } from "react-router-dom";
import "./css/custom.css";
import "./App.css";

// Render it in your app!
function App(props) {
  const [selected, setSelection] = useState(false);
  let colorMode = "light";
  if (selected) {
    colorMode = "dark";
  } else {
    colorMode = "light";
  }
  return (
    <Provider theme={defaultTheme} height="100%" colorScheme={colorMode}>
      <Flex
        direction="column"
        width="calc(100%)"
        gap="size-100"
        borderWidth="thin"
        borderColor="dark"
        height="calc(100%)"
      >
        <View backgroundColor="gray-200" height="size-400">
          <AppHeader setSelection={setSelection} />
        </View>
        <View
          backgroundColor="gray-50"
          width="calc(max(80%, size-6000)"
          height="100%"
          alignSelf="center"
          flex="true"
          borderWidth="thin"
          borderColor="dark"
          borderRadius="medium"
          padding="size-100"
        >
          <Outlet />
        </View>
        <View backgroundColor="gray-200" height="size-400">
          <AppFooter />
        </View>
      </Flex>
    </Provider>
  );
}

export default App;
