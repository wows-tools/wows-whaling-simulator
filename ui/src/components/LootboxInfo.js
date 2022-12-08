import React from 'react';
import axios from 'axios';
import {Image} from '@adobe/react-spectrum';
import {Link} from '@adobe/react-spectrum';
import {Text} from '@adobe/react-spectrum';
import {Heading} from '@adobe/react-spectrum';
import {View} from '@adobe/react-spectrum';
import {Flex} from '@adobe/react-spectrum';
import {Divider} from '@adobe/react-spectrum';
import {IllustratedMessage} from '@adobe/react-spectrum';
import {ActionButton, Dialog, ButtonGroup, Button, DialogTrigger, Slider, Picker, Item, SearchField} from '@adobe/react-spectrum';
import {Content} from '@adobe/react-spectrum';
import {Link as RouterLink} from 'react-router-dom';
import {useParams} from 'react-router-dom';
import NotFound from '@spectrum-icons/illustrations/NotFound';
import Money from '@spectrum-icons/workflow/Money';
import User from '@spectrum-icons/workflow/User';

import { API_ROOT } from '../api-config';

function withParams(Component) {
  return props => <Component {...props} params={useParams()} />;
}

class LootboxInfo extends React.Component {
  	state = {
    		lootbox: null
  	}

  componentDidMount() {
    let { lootboxId } = this.props.params;
    axios.get(`${API_ROOT}/api/v1/lootboxes/${lootboxId}`)
      .then(res => {
        const lootbox = res.data;
	console.log(res.data);
        this.setState({ lootbox });
      })
  }

  render() { 
	let lootboxContent;
	if (this.state.lootbox) {
		lootboxContent = (
        	<IllustratedMessage>
        	<Image height="200px" objectFit="scale-down" src={API_ROOT + this.state.lootbox.img} alt={this.state.lootbox.name} />
        	<Content>{this.state.lootbox.name}</Content>
        	</IllustratedMessage>

		)
	} else {
		lootboxContent = (
    		<IllustratedMessage>
    		  <NotFound />
    		  <Heading>No result</Heading>
    		  <Content>Container found</Content>
    		</IllustratedMessage>
		)
	}
    return (
      <Flex alignContent="center" justifyContent="center" margin="size-100" direction="column" gap="size-100" justifyContent="center" alignContent="center" alignItems="center">
      <View backgroundColor="gray-200" borderColor="dark" borderRadius="small" width="size-2400">
        {lootboxContent}
      </View>
      <View>
      <DialogTrigger>
      <ActionButton><Money /><Text>Start Whaling!</Text></ActionButton>
      {(close) => (
        <Dialog>
          <Heading>Whaling Parameters</Heading>
          <Divider />
          <Content>Please, fill the whaling parameters:</Content>
	  <Slider label="Number of Containers to open" defaultValue="12" maxValue="1000" />
	  <Picker label="Realm/Wows Server">
  <Item key="na">NA</Item>
  <Item key="eu">EU</Item>
  <Item key="asia">Asia</Item>
</Picker>
	  <SearchField label="Search for users" icon={<User />} />
          <ButtonGroup>
            <Button variant="secondary" onPress={close}>Cancel</Button>
            <Button variant="accent" onPress={close} autoFocus>Whale</Button>
          </ButtonGroup>
        </Dialog>
      )}
      </DialogTrigger>

      </View>
      <Divider size="M" />
	
      </Flex>
    )
  }
}

export default withParams(LootboxInfo);
