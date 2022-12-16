import React from "react";
import { useEffect } from "react";
import axios from "axios";
import { Image } from "@adobe/react-spectrum";
import LootboxContent from "./LootboxContent";
import { Link } from "@adobe/react-spectrum";
import { Text } from "@adobe/react-spectrum";
import { Heading } from "@adobe/react-spectrum";
import { View } from "@adobe/react-spectrum";
import { Flex } from "@adobe/react-spectrum";
import { ContextualHelp } from "@adobe/react-spectrum";
import { Form } from "@adobe/react-spectrum";
import { Switch } from "@adobe/react-spectrum";
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
import Back from "@spectrum-icons/workflow/Back";
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

function StatsWhalingResult(props) {
  // TODO
  return <></>;
}

function SimpleWhalingResult(props) {
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

function WhalingResults(props) {
  switch (props.whalingData.simulation_type) {
    case "simple_whaling_quantity":
    case "simple_whaling_target":
      return <SimpleWhalingResult whalingData={props.whalingData} />;
    case "stats_whaling_quantity":
    case "stats_whaling_target":
      return <StatsWhalingResult whalingData={props.whalingData} />;
  }
}

export default WhalingResults;
