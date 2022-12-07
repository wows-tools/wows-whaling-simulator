import React from 'react';
import axios from 'axios';
import {Image} from '@adobe/react-spectrum';
import {Link} from '@adobe/react-spectrum';
import {Link as RouterLink} from 'react-router-dom';

import { API_ROOT } from '../api-config';


export default class LootboxList extends React.Component {
  state = {
    lootboxes: []
  }

  componentDidMount() {
    axios.get(`${API_ROOT}/api/v1/lootboxes`)
      .then(res => {
        const lootboxes = res.data['lootboxes'];
	console.log(res.data);
        this.setState({ lootboxes });
      })
  }

  logMapElements(value, key, map) {
    console.log(`m[${key}] = ${value}`);
  }

  render() {
    return (
      <ul>
	{console.log(this.state.lootboxes)}
        {
          this.state.lootboxes.map(lb => 
		  <Link>
  		  <RouterLink to={"/lootboxes/" + lb.id}>
		  <Image height="200px" objectFit="scale-down" src={API_ROOT + lb.img} alt={lb.name} />
		  </RouterLink>
		  </Link>
	  )
        }
      </ul>
    )
  }
}

