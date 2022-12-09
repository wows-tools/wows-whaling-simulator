import React from 'react';
import {useEffect} from 'react';
import axios from 'axios';
import {Image} from '@adobe/react-spectrum';
import {Link} from '@adobe/react-spectrum';
import {Text} from '@adobe/react-spectrum';
import {Heading} from '@adobe/react-spectrum';
import {View} from '@adobe/react-spectrum';
import {Flex} from '@adobe/react-spectrum';
import {Form} from '@adobe/react-spectrum';
import {Divider} from '@adobe/react-spectrum';
import {IllustratedMessage} from '@adobe/react-spectrum';
import {ComboBox, ActionButton, AlertDialog, ButtonGroup, Button, DialogTrigger, Slider, Picker, Item, SearchField, DialogContainer, TextField} from '@adobe/react-spectrum';
import {Content} from '@adobe/react-spectrum';
import {Link as RouterLink} from 'react-router-dom';
import {useParams} from 'react-router-dom';
import NotFound from '@spectrum-icons/illustrations/NotFound';
import Money from '@spectrum-icons/workflow/Money';
import User from '@spectrum-icons/workflow/User';
import {useAsyncList} from 'react-stately';

import { API_ROOT } from '../api-config';

function checkUnset(props) {
        return ((props === undefined) || props === null || (props.length === 0))
}

function WhaleBox(props) {
  const [isOpen, setOpen] = React.useState(false);
  const [whaling, setWhaling] = React.useState(false);
  const [realm, setRealm] = React.useState();
  const [player, setPlayer] = React.useState();
  const [numlootbox, setNumlootbox] = React.useState(25);

  let realmOptions = [
    { id: "eu", name: 'EU' },
    { id: "na", name: 'NA' },
    { id: "asia", name: 'Asia' },
  ];

  let list = useAsyncList({
    async load({ signal, cursor, filterText }) {
      if (filterText.length < 3) {
	return {
	  items: []
	}
      }
      let res = await fetch(
	`${API_ROOT}/api/v1/realms/${realm}/players?nick_start=${filterText}`,
	{ signal }
      );
      let json = await res.json();

      return {
	items: json.players
      };
    }
  });

  const triggerWhaling = () => {
    axios.get(`${API_ROOT}/api/v1/lootboxes/${props.lootboxId}/realms/${realm}/players/${player}/simple_whaling_quantity?number_lootbox=${numlootbox}`) 
      .then(res => {
	const stats = res.data;
	props.setStats(stats);
      })
      setPlayer("")
      list.setFilterText("")
  }

  return (
    <>
    <ActionButton onPress={() => setOpen(true)}>	
    <Money /><Text>Start Whaling!</Text>
    </ActionButton>
    <DialogContainer onDismiss={() => setOpen(false)}>
    {isOpen &&
      <AlertDialog
      title="Let's Do Some Whaling"
      variant="confirmation"
      primaryActionLabel="Start Whaling"
      cancelLabel="Cancel"
      onPrimaryAction={triggerWhaling}
      isPrimaryActionDisabled={checkUnset(player)}
      >

      <Form maxWidth="size-3600">
      <Picker
      label="Realm/Wows Server"
      items={realmOptions}
      onSelectionChange={(selected) => setRealm(selected)}
      autoFocus="true"
      >
      {(item) => <Item>{item.name}</Item>}
      </Picker>
      <ComboBox
      label="Player Search"
      items={list.items}
      inputValue={list.filterText}
      onInputChange={list.setFilterText}
      loadingState={list.loadingState}
      isDisabled={checkUnset(realm)}
      onSelectionChange={(selected) => setPlayer(selected)}
      >
      {(item) => <Item key={item.account_id}>{item.nickname}</Item>}
      </ComboBox>

      <Slider label="Containers Quantity" defaultValue="25" maxValue="1000" onChange={setNumlootbox}/>
      </Form>

      </AlertDialog>
    }
    </DialogContainer>
    </>
  )
}


function LootboxInfo() {
  const [stats, setStats] = React.useState(false);
  const [lootbox, setLootbox] = React.useState(false);
  let { lootboxId } = useParams();

  const [isOpen, setOpen] = React.useState(false);
  useEffect(() => {
    axios.get(`${API_ROOT}/api/v1/lootboxes/${lootboxId}`)
      .then(res => {
	const lootbox = res.data;
	setLootbox(lootbox);
      })
  },[lootboxId]);

  const validateInput = () => {
    return false
  }
  let lootboxContent;
  if (lootbox) {
    lootboxContent = (
      <IllustratedMessage>
      <Image height="200px" objectFit="scale-down" src={API_ROOT + lootbox.img} alt={lootbox.name} />
      <Content>{lootbox.name}</Content>
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
    <Flex margin="size-100" direction="column" gap="size-100" justifyContent="center" alignContent="center" alignItems="center">
    <View backgroundColor="gray-200" borderColor="dark" borderRadius="small" width="size-2400">
    {lootboxContent}
    </View>
    <View>
    <WhaleBox lootboxId={lootboxId} setStats={setStats} />
    </View>

    <Divider size="M" />


    </Flex>
  )
}

export default LootboxInfo;
