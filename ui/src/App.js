import {Flex, View, Provider, defaultTheme} from '@adobe/react-spectrum';
import Footer from './components/Footer';
import LootboxList from './components/LootboxList';
import { Outlet } from "react-router-dom";
import './css/custom.css'
import './App.css'

// Render it in your app!
function App(props) {
  return (
    <Provider theme={defaultTheme} height="100%" colorScheme="dark">
	<Flex direction="column" width="calc(100%)" gap="size-100" borderWidth="thin" borderColor="dark" height="calc(100% - size-600)">
  		<View backgroundColor="gray-200" height="size-600">
	  	<h2>WoWs Whaling Simulator (version alpha.alpha.alpha)</h2>
	        </View>
  		<View backgroundColor="gray-400"  width="calc(max(80%, size-6000)" 
		  height="100%" alignSelf="center" marging="size-100" flex="true" borderWidth="thin" borderColor="dark" borderRadius="medium">
	        <Outlet />
	        </View>
	</Flex>
	<Footer/>
    </Provider>
  );
}

export default App;
