import React from 'react';
import {useEffect} from 'react';
import axios from 'axios';
import {Image} from '@adobe/react-spectrum';
import {Link} from '@adobe/react-spectrum';
import {Text} from '@adobe/react-spectrum';
import {Heading} from '@adobe/react-spectrum';
import {View} from '@adobe/react-spectrum';
import {Flex} from '@adobe/react-spectrum';
import {ContextualHelp} from '@adobe/react-spectrum';
import {Form} from '@adobe/react-spectrum';
import {Divider} from '@adobe/react-spectrum';
import {IllustratedMessage} from '@adobe/react-spectrum';
import {ComboBox, ActionButton, AlertDialog, ButtonGroup, Button, DialogTrigger, Slider, Picker, Item, SearchField, DialogContainer, TextField} from '@adobe/react-spectrum';
import {Content} from '@adobe/react-spectrum';
import {TableView, TableHeader, Column, TableBody, Row, Cell} from '@adobe/react-spectrum';
import {ListBox} from '@adobe/react-spectrum';
import {Section} from '@adobe/react-spectrum';
import {Link as RouterLink} from 'react-router-dom';
import {useParams} from 'react-router-dom';
import NotFound from '@spectrum-icons/illustrations/NotFound';
import Money from '@spectrum-icons/workflow/Money';
import User from '@spectrum-icons/workflow/User';
import Star from '@spectrum-icons/workflow/Star';
import {useAsyncList} from 'react-stately';

import { API_ROOT } from '../api-config';

function checkUnset(props) {
  return ((props === undefined) || props === null || (props.length === 0))
}

function ShipInfo() {
	return (
		<ContextualHelp variant="info" placement="top start">
		  <Content>
		    <Text>
			'*' indicates rare ships
			(not obtanable for resources or money)
		    </Text>
		  </Content>
		</ContextualHelp>
	)
}

function RenderShipList(props) {
  var groupSize = 5;
  var rows = props.ships.map(function(ship) {
    // map content to Item
    if ('rare' in ship.attributes) {
      return (<Item><Text>
	{ship.name} *
	</Text></Item>);
    } else {
      return (<Item>{ship.name}</Item>);
    }
  }).reduce(function(r, element, index) {
    // create element groups with size 5, result looks like:
    // [[elem1, elem2, elem3], [elem4, elem5, elem6], ...]
    index % groupSize === 0 && r.push([]);
    r[r.length - 1].push(element);
    return r;
  }, []).map(function(rowContent) {
    // surround every group with 'row'
    return (<ListBox width="size-2400" selectionMode="none" aria-label="Alignment">{rowContent}</ListBox>);
  });

  return (
    <Flex direction="row" gap="size-100" wrap>
    {rows.map((row) => (
      <View>
      {row}
      </View>
    ))}
    </Flex>
  )
}


function RenderItems(props) {
  return (
    <TableView selectionMode="none" density="compact">
    <TableHeader>
    <Column key="name" width="60%">Name</Column>
    <Column key="quantity" width="40%">Quantity</Column>
    </TableHeader>
    <TableBody>
    {props.items.map((item) =>
    <Row>
      <Cell>{item.name}</Cell>
      <Cell>{item.quantity}</Cell>
    </Row>
    )}
    </TableBody>
    </TableView>)
}

function WhalingResult(props) {
  console.log(props.whalingData)
  let ship_cat = {tx: [], tix_viii: [], tvii_: []}
  for (const ship of props.whalingData.collectables_items) {
    switch (ship.attributes.tier) {
      case 'X': ship_cat.tx.push(ship);break;
      case 'IX':
      case 'VIII':
	ship_cat.tix_viii.push(ship);break;
      default:
	ship_cat.tvii_.push(ship);break;
    } 
  }
  let other_cat = {resource: [], eco: [], other: []}
  for (const item of props.whalingData.other_items) {
    switch (item.attributes.type) {
      case 'economic bonus': other_cat.eco.push(item);break;
      case 'resource': other_cat.resource.push(item);break;
      default: other_cat.other.push(item);break;
    } 
  }

  return (
    <Flex direction="column" gap="size-100">
    <View>
      <Flex direction="row" gap="size-100">
      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Tier X ships <ShipInfo/></Heading>
      <Divider size="M" />
      <RenderShipList ships={ship_cat.tx}/>
      </View>

      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Tier IX & VIII ships <ShipInfo/></Heading>
      <Divider size="M" />
      <RenderShipList ships={ship_cat.tix_viii}/>
      </View>

      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Tier VII & bellow ships <ShipInfo/></Heading>
      <Divider size="M" />
      <RenderShipList ships={ship_cat.tvii_}/>
      </View>
      </Flex>
    </View>
    <View>
      <Flex direction="row" gap="size-100">
      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Resource</Heading>
      <RenderItems items={other_cat.resource}/>
      </View>

      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Economic Bonus</Heading>
      <RenderItems items={other_cat.eco}/>
      </View>

      <View width="33%"  borderRadius="medium" borderWidth="thin" borderColor="dark" padding="size-100" overflow="scroll" maxHeight="size-5000">
      <Heading>Other</Heading>
      <RenderItems items={other_cat.other}/>
      </View>
      </Flex>
    </View>
    </Flex>
  );
}

function WhaleBox(props) {
  const [isOpen, setOpen] = React.useState(false);
  const [whaling, setWhaling] = React.useState(false);
  const [realm, setRealm] = React.useState();
  const [player, setPlayer] = React.useState();
  const [numlootbox, setNumlootbox] = React.useState(20);

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
    setPlayer()
    setRealm()
    list.setFilterText("")
  }

  return (
    <>
    <Button variant="accent" style="fill" onPress={() => setOpen(true)}>	
    <Money /><Text>Start Whaling!</Text>
    </Button>
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

      <Form maxWidth="size-4600">
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

      <Slider height="size-1000" width="size-3600" label="Containers Quantity" defaultValue="20" maxValue="1000" onChange={setNumlootbox}/>
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
    { stats && (
      <View>
      <Heading>Whaling results:</Heading>
      <WhalingResult whalingData={stats}/>
      </View>
    )}


    </Flex>
  )
}

export default LootboxInfo;
