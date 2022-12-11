import React from "react";
import { useEffect } from "react";
import axios from "axios";
import { Image } from "@adobe/react-spectrum";
import { Link } from "@adobe/react-spectrum";
import { Text } from "@adobe/react-spectrum";
import { Heading } from "@adobe/react-spectrum";
import { View } from "@adobe/react-spectrum";
import { Flex } from "@adobe/react-spectrum";
import { ContextualHelp } from "@adobe/react-spectrum";
import { Form } from "@adobe/react-spectrum";
import { Grid } from "@adobe/react-spectrum";
import { Divider } from "@adobe/react-spectrum";
import { IllustratedMessage } from "@adobe/react-spectrum";
import { NumberField } from "@adobe/react-spectrum";
import {
  ComboBox,
  ActionButton,
  AlertDialog,
  ButtonGroup,
  Button,
  DialogTrigger,
  Slider,
  Picker,
  Item,
  SearchField,
  DialogContainer,
  TextField,
} from "@adobe/react-spectrum";
import { Content } from "@adobe/react-spectrum";
import {
  Tabs,
  TabList,
  TabPanels,
  TableView,
  TableHeader,
  Column,
  TableBody,
  Row,
  Cell,
} from "@adobe/react-spectrum";
import { ListBox } from "@adobe/react-spectrum";
import { useNavigate } from "react-router-dom";
import { Section } from "@adobe/react-spectrum";
import { Link as RouterLink } from "react-router-dom";
import { useParams } from "react-router-dom";
import NotFound from "@spectrum-icons/illustrations/NotFound";
import Money from "@spectrum-icons/workflow/Money";
import User from "@spectrum-icons/workflow/User";
import Star from "@spectrum-icons/workflow/Star";
import { useAsyncList } from "react-stately";
import GenericTile from "./GenericTile";

import { API_ROOT } from "../api-config";

function checkUnset(props) {
  return props === undefined || props === null || props.length === 0;
}

function ShipInfo() {
  return (
    <ContextualHelp variant="info" placement="top start">
      <Content>
        <Text>
          '*' indicates rare ships (not obtanable for resources or money)
        </Text>
      </Content>
    </ContextualHelp>
  );
}

function RenderShipList(props) {
  var groupSize = 5;
  var rows = props.ships
    .map(function (ship) {
      // map content to Item
      if ("rare" in ship.attributes) {
        return (
          <Item>
            <Text>{ship.name} *</Text>
          </Item>
        );
      } else {
        return <Item>{ship.name}</Item>;
      }
    })
    .reduce(function (r, element, index) {
      // create element groups with size 5, result looks like:
      // [[elem1, elem2, elem3], [elem4, elem5, elem6], ...]
      index % groupSize === 0 && r.push([]);
      r[r.length - 1].push(element);
      return r;
    }, [])
    .map(function (rowContent) {
      // surround every group with 'row'
      return (
        <ListBox width="size-2400" selectionMode="none">
          {rowContent}
        </ListBox>
      );
    });

  return (
    <View
      width="33%"
      borderRadius="medium"
      borderWidth="thin"
      borderColor="dark"
      padding="size-100"
      overflow="auto"
      backgroundColor="gray-100"
      maxHeight="size-5000"
    >
      <Heading>
        {props.title}
        <ShipInfo />
      </Heading>
      <Divider size="M" />

      <Flex direction="row" gap="size-100" wrap>
        {rows.map((row) => (
          <View>{row}</View>
        ))}
      </Flex>
    </View>
  );
}

function RenderItems(props) {
  return (
    <View
      width="33%"
      borderRadius="medium"
      borderWidth="thin"
      borderColor="dark"
      padding="size-100"
      backgroundColor="gray-100"
      overflow="auto"
      maxHeight="size-5000"
    >
      <Heading>{props.title}</Heading>
      <Divider size="M" />

      <TableView selectionMode="none" density="compact">
        <TableHeader>
          <Column key="name" width="60%">
            Name
          </Column>
          <Column key="quantity" width="40%">
            Quantity
          </Column>
        </TableHeader>
        <TableBody>
          {props.items.map((item) => (
            <Row>
              <Cell>{item.name}</Cell>
              <Cell>{item.quantity}</Cell>
            </Row>
          ))}
        </TableBody>
      </TableView>
    </View>
  );
}

function WhalingResult(props) {
  let ship_cat = { tx: [], tix_viii: [], tvii_: [] };
  for (const ship of props.whalingData.collectables_items) {
    switch (ship.attributes.tier) {
      case "X":
        ship_cat.tx.push(ship);
        break;
      case "IX":
      case "VIII":
        ship_cat.tix_viii.push(ship);
        break;
      default:
        ship_cat.tvii_.push(ship);
        break;
    }
  }
  let other_cat = { resource: [], eco: [], other: [] };
  for (const item of props.whalingData.other_items) {
    switch (item.attributes.type) {
      case "economic bonus":
        other_cat.eco.push(item);
        break;
      case "resource":
        other_cat.resource.push(item);
        break;
      default:
        other_cat.other.push(item);
        break;
    }
  }

  return (
    <Flex direction="column" gap="size-100">
      <Grid
        areas={["slot1 slot2 slot3"]}
        gap="size-100"
        justifyItems="center"
        wrap
      >
        <View width="size-3600">
          <GenericTile
            header="Doubloons"
            subheader="Doubloons (and real money spent)"
            scale="Doubloons"
            number={props.whalingData.game_money_spent}
            footer={
              "(ie: €" +
              props.whalingData.euro_spent +
              " or $" +
              props.whalingData.dollar_spent +
              ")"
            }
            minWidth="size-3600"
          />
        </View>
        <View width="size-3600">
          <GenericTile
            header="Opened"
            subheader="Container Opened"
            scale="Container(s)"
            number={props.whalingData.container_opened}
            minWidth="size-3600"
          />
        </View>
        <View width="size-3600">
          <GenericTile
            header="Pities"
            subheader="Pity trigger count"
            scale="Pities"
            number={props.whalingData.pities}
            minWidth="size-3600"
          />
        </View>
      </Grid>
      <View>
        <Flex direction="row" gap="size-100">
          <RenderShipList ships={ship_cat.tx} title="Tier X" />
          <RenderShipList ships={ship_cat.tix_viii} title="Tier IX & VIII" />
          <RenderShipList ships={ship_cat.tvii_} title="Tier VII & bellow" />
        </Flex>
      </View>
      <View>
        <Flex direction="row" gap="size-100">
          <RenderItems items={other_cat.resource} title="Resources" />
          <RenderItems items={other_cat.eco} title="Economic Bonuses" />
          <RenderItems items={other_cat.other} title="Oter Items" />
        </Flex>
      </View>
    </Flex>
  );
}

function WhaleBox(props) {
  const [isOpen, setOpen] = React.useState(false);
  const [realm, setRealm] = React.useState();
  const [player, setPlayer] = React.useState();
  const [numlootbox, setNumlootbox] = React.useState(20);
  const navigate = useNavigate();

  let realmOptions = [
    { id: "eu", name: "EU" },
    { id: "na", name: "NA" },
    { id: "asia", name: "Asia" },
  ];

  const GoHome = () => {
    navigate("/");
  };

  const setRealmReset = (value) => {
    // Reset Player when changing Realm
    setPlayer();
    list.setFilterText("");
    setRealm(value);
  };

  let list = useAsyncList({
    async load({ signal, cursor, filterText }) {
      if (filterText.length < 3) {
        return {
          items: [],
        };
      }
      let res = await fetch(
        `${API_ROOT}/api/v1/realms/${realm}/players?nick_start=${filterText}`,
        { signal }
      );
      let json = await res.json();

      return {
        items: json.players,
      };
    },
  });

  const triggerWhaling = () => {
    axios
      .get(
        `${API_ROOT}/api/v1/lootboxes/${props.lootboxId}/realms/${realm}/players/${player}/simple_whaling_quantity?number_lootbox=${numlootbox}`
      )
      .then((res) => {
        const stats = res.data;
        props.setStats(stats);
        props.setTab("whaling");
      });
  };

  return (
    <>
      <Flex direction="row" gap="size-400">
        <Button variant="accent" style="fill" onPress={() => setOpen(true)}>
          <Money />
          <Text>Start Whaling!</Text>
        </Button>
        <Button variant="secondary" style="outline" onPress={GoHome}>
          <Text>Back to Containers</Text>
        </Button>
      </Flex>

      <DialogContainer onDismiss={() => setOpen(false)}>
        {isOpen && (
          <AlertDialog
            title="Let's Do Some Whaling"
            variant="confirmation"
            primaryActionLabel="Start Whaling"
            cancelLabel="Cancel"
            onPrimaryAction={triggerWhaling}
            isPrimaryActionDisabled={checkUnset(player)}
            minWidth="size-6000"
          >
            <Form>
              <Picker
                label="Realm/Wows Server"
                items={realmOptions}
                onSelectionChange={(selected) => setRealmReset(selected)}
                autoFocus="true"
                defaultSelectedKey={realm}
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

              <Flex direction="row" gap="size-100">
                <View>
                  <Slider
                    height="size-1000"
                    label="Number of containers"
                    value={numlootbox}
                    width="size-3600"
                    maxValue="1000"
                    showValueLabel={false}
                    onChange={setNumlootbox}
                  />
                </View>
                <View marginTop="calc(single-line-height / 2)">
                  <NumberField
                    width="size-1200"
                    value={numlootbox}
                    minValue={0}
                    maxValue="1000"
                    onChange={setNumlootbox}
                  />
                </View>
              </Flex>
            </Form>
          </AlertDialog>
        )}
      </DialogContainer>
    </>
  );
}

function RenderSlot(props) {
  // FIXME display properly the items
  return (
    <TableView selectionMode="none" density="compact" overflowMode="wrap">
      <TableHeader>
        <Column key="name">Category</Column>
        <Column key="droprate">Drop Rate</Column>
        <Column key="item_pool_size">Item Pool Size</Column>
        <Column key="item_drop_rate">Individual Item Drop Rate</Column>
        <Column key="items">Items</Column>
      </TableHeader>
      <TableBody>
        {Object.values(props.drops).map((cat, index) => {
          return (
            <Row>
              <Cell>{cat.name}</Cell>
              <Cell>{cat.drop_rate} %</Cell>
              <Cell>{cat.items.length}</Cell>
              <Cell>{cat.drop_rate / cat.items.length} %</Cell>
              <Cell>
                <View maxHeight="size-2000" overflow="auto">
                  <ul>
                    {cat.items.map((item) => (
                      <li>
                        {item.name} (x {item.quantity})
                      </li>
                    ))}
                  </ul>
                </View>
              </Cell>
            </Row>
          );
        })}
      </TableBody>
    </TableView>
  );
}

function LootboxContent(props) {
  let lootbox = props.lootbox;

  if (!props.lootbox) {
    return (
      <IllustratedMessage>
        <NotFound />
        <Heading>No result</Heading>
        <Content>Container found</Content>
      </IllustratedMessage>
    );
  }

  return (
    <View>
      <Grid
        areas={["slot1 slot2 slot3 slot4"]}
        gap="size-100"
        justifyItems="center"
        wrap
      >
        <View
          width="size-3600"
          backgroundColor="gray-100"
          borderRadius="medium"
          borderWidth="thin"
          borderColor="dark"
          padding="size-100"
        >
          <IllustratedMessage>
            <Image
              height="size-2000"
              objectFit="scale-down"
              src={API_ROOT + lootbox.img}
              alt={lootbox.name}
            />
            <Content>{lootbox.name}</Content>
          </IllustratedMessage>
        </View>
        <View width="size-3600">
          <GenericTile
            header="Price"
            subheader="Container Price"
            scale="Doubloons"
            number={lootbox.price}
            minWidth="size-3600"
          />
        </View>
        <View width="size-3600">
          <GenericTile
            header="Pity"
            subheader="Pity count"
            scale="Containers"
            number={lootbox.pity}
            minWidth="size-3600"
          />
        </View>
        <View width="size-3600">
          <GenericTile
            header="Slots"
            subheader="number of slots"
            scale="Slot(s)"
            number={lootbox.drops.length}
            minWidth="size-3600"
          />
        </View>
      </Grid>
      <Tabs>
        <TabList>
          {props.lootbox.drops.map((drops, index) => (
            <Item key={index + 1}>Slot {index + 1}</Item>
          ))}
        </TabList>
        <TabPanels>
          {props.lootbox.drops.map((drops, index) => (
            <Item key={index + 1}>
              <RenderSlot drops={drops} />
            </Item>
          ))}
        </TabPanels>
      </Tabs>
    </View>
  );
}

function Lootbox() {
  const [stats, setStats] = React.useState(false);
  const [tabSelected, setTabSelected] = React.useState("container");
  const [lootbox, setLootbox] = React.useState(false);
  let { lootboxId } = useParams();

  const [isOpen, setOpen] = React.useState(false);
  useEffect(() => {
    axios.get(`${API_ROOT}/api/v1/lootboxes/${lootboxId}`).then((res) => {
      const lootbox = res.data;
      setLootbox(lootbox);
    });
  }, [lootboxId]);

  const validateInput = () => {
    return false;
  };
  let defaultTab = "container";
  let disabledTabs = ["whaling"];
  if (stats) {
    defaultTab = "whaling";
    disabledTabs = [];
  }
  return (
    <Flex
      margin="size-100"
      direction="column"
      gap="size-100"
      justifyContent="center"
      alignContent="center"
      alignItems="center"
    >
      <View>
        <WhaleBox
          lootboxId={lootboxId}
          setStats={setStats}
          setTab={setTabSelected}
        />
      </View>

      <Divider size="S" />
      <Tabs
        disabledKeys={disabledTabs}
        selectedKey={tabSelected}
        onSelectionChange={setTabSelected}
      >
        <TabList>
          <Item key="container">Container Drop Rates</Item>
          <Item key="whaling">Whaling Session</Item>
        </TabList>
        <TabPanels>
          <Item key="container">
            <Heading>Lootbox Info:</Heading>
            <LootboxContent lootbox={lootbox} />
          </Item>
          {stats && (
            <Item key="whaling">
              <View>
                <Heading>Whaling results:</Heading>
                <WhalingResult whalingData={stats} />
              </View>
            </Item>
          )}
        </TabPanels>
      </Tabs>
    </Flex>
  );
}

export default Lootbox;
