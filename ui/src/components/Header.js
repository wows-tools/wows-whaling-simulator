import {
  Flex,
  Header,
  View,
  Provider,
  defaultTheme,
  Switch,
  Text,
} from "@adobe/react-spectrum";

function AppHeader(props) {
  return (
    <Header>
      <Flex direction="row">
        <Flex width="50%" alignContent="left" justifyContent="left"></Flex>
        <Flex width="50%" alignContent="right" justifyContent="right">
          <Switch onChange={props.setSelection}>Switch Dark Mode</Switch>
        </Flex>
      </Flex>
    </Header>
  );
}

export default AppHeader;
