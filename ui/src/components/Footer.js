import "../css/custom.css";
import { Content, Flex, View, Footer, Link } from "@adobe/react-spectrum";

// Render it in your app!
function AppFooter() {
  return (
    <Footer>
      <Flex alignContent="center" justifyContent="center">
        <View backgroundColor="gray-200" height="size-100">
          <Content>
            <a href="https://github.com/kakwa/wows-whaling-simulator">
              WoWs Whaling Simulator (Source Code)
            </a>{" "}
            • version alpha.alpha • © 2022 • Kakwa • Released under the MIT
            Public License
          </Content>
        </View>
      </Flex>
    </Footer>
  );
}

export default AppFooter;
