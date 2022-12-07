import React from 'react';
import axios from 'axios';
import {Image} from '@adobe/react-spectrum';
import {Link} from '@adobe/react-spectrum';
import {Text} from '@adobe/react-spectrum';
import {Heading} from '@adobe/react-spectrum';
import {View} from '@adobe/react-spectrum';
import {IllustratedMessage} from '@adobe/react-spectrum';
import {Content} from '@adobe/react-spectrum';
import {Link as RouterLink} from 'react-router-dom';
import {useParams} from 'react-router-dom';
import NotFound from '@spectrum-icons/illustrations/NotFound';

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
      <View>
        {lootboxContent}
      </View>
    )
  }
}

export default withParams(LootboxInfo);
