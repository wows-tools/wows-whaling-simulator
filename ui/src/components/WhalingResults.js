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
      if ("rare" in ship.attributes && ship.attributes["rare"] === "true") {
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
      minWidth="size-3600"
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
      minWidth="size-3600"
      borderRadius="medium"
      borderWidth="thin"
      borderColor="dark"
      padding="size-100"
      backgroundColor="gray-100"
      overflow="auto"
      maxHeight="size-5000"
    >
      <Heading>
        {props.prefix}
        {props.title}
      </Heading>
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
      <Flex direction="row" gap="size-100" justifyContent="space-evenly" wrap>
        <View width="size-3600">
          <GenericTile
            header="Doubloons"
            subheader="Doubloons (and real money) spent"
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
            header="Ships"
            subheader="Number of ships dropped"
            scale="Ships"
            number={props.whalingData.collectables_items.length}
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
      </Flex>
      <Divider size="S" />
      <View>
        <Flex direction="row" gap="size-100" wrap>
          <RenderShipList ships={ship_cat.tx} title="Tier X" />
          <RenderShipList ships={ship_cat.tix_viii} title="Tier IX & VIII" />
          <RenderShipList ships={ship_cat.tvii_} title="Tier VII & bellow" />
        </Flex>
      </View>
      <Divider size="S" />
      <View>
        <Flex direction="row" gap="size-100" wrap>
          <RenderItems items={other_cat.resource} title="Resources" />
          <RenderItems items={other_cat.eco} title="Economic Bonuses" />
          <RenderItems items={other_cat.other} title="Other Items" />
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

function RenderPercentiles(props) {
  return (
    <View
      width="33%"
      minWidth="size-3600"
      borderRadius="medium"
      borderWidth="thin"
      borderColor="dark"
      padding="size-100"
      backgroundColor="gray-100"
      overflow="auto"
      maxHeight="size-5000"
    >
      <Heading>
        Containers required to have X% chance of getting "{props.target}"
      </Heading>
      <Divider size="M" />

      <TableView selectionMode="none" density="compact">
        <TableHeader>
          <Column key="name">Odds</Column>
          <Column key="quantity">Number of Containers</Column>
        </TableHeader>
        <TableBody>
          {Object.keys(props.percentiles).map((key) => (
            <Row>
              <Cell>{(isNaN(key) && <>{key}</>) || <>{key}%</>}</Cell>
              <Cell>{props.percentiles[key]}</Cell>
            </Row>
          ))}
        </TableBody>
      </TableView>
    </View>
  );
}

function RenderByAttribute(props) {
  let attribute = props.attribute;
  let data = props.data;
  return (
    <View
      width="33%"
      minWidth="size-3600"
      borderRadius="medium"
      borderWidth="thin"
      borderColor="dark"
      padding="size-100"
      backgroundColor="gray-100"
      overflow="auto"
      maxHeight="size-5000"
    >
      <Heading>{attribute}</Heading>
      <Divider size="M" />

      <TableView selectionMode="none" density="compact">
        <TableHeader>
          <Column key="name">Name</Column>
          <Column key="quantity">Quantity</Column>
        </TableHeader>
        <TableBody>
          {Object.keys(data).map((key) => (
            <Row>
              <Cell>{key}</Cell>
              <Cell>{data[key]}</Cell>
            </Row>
          ))}
        </TableBody>
      </TableView>
    </View>
  );
}

function StatsWhalingResult(props) {
  let other_cat = { resource: [], eco: [], other: [] };
  for (const [key, item] of Object.entries(props.whalingData.avg_by_item)) {
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
  if (!("tier" in props.whalingData.avg_by_attribute)) {
    props.whalingData.avg_by_attribute["tier"] = {};
  }
  if (!("rare" in props.whalingData.avg_by_attribute)) {
    props.whalingData.avg_by_attribute["rare"] = {};
  }

  return (
    <Flex direction="column" gap="size-100">
      <Flex direction="row" gap="size-100" justifyContent="space-evenly" wrap>
        <View width="size-4200">
          <GenericTile
            header="Simulation runs"
            subheader="Number of runs"
            scale="Runs"
            number={props.whalingData.session_count}
            footer={
              props.whalingData.total_opened + " containers opened in total"
            }
            minWidth="size-4200"
          />
        </View>
        <View width="size-4200">
          <GenericTile
            header="Average Doubloons"
            subheader="Doubloons (and real money) spent"
            scale="Doubloons"
            number={props.whalingData.avg_game_money_spent}
            footer={
              "(ie: €" +
              props.whalingData.avg_euro_spent +
              " or $" +
              props.whalingData.avg_dollar_spent +
              ")"
            }
            minWidth="size-4200"
          />
        </View>
        <View width="size-4200">
          <GenericTile
            header="Average Opened"
            subheader="Container Opened"
            scale="Container(s)"
            number={props.whalingData.avg_opened}
            minWidth="size-4200"
          />
        </View>
        <View width="size-4200">
          <GenericTile
            header="Average Pities"
            subheader="Pity trigger count"
            scale="Pities"
            number={props.whalingData.avg_pities}
            minWidth="size-4200"
          />
        </View>
      </Flex>

      <Divider size="S" />
      <View>
        <Flex direction="row" gap="size-100" justifyContent="space-evenly" wrap>
          {props.whalingData.simulation_type == "stats_whaling_target" && (
            <RenderPercentiles
              target={props.whalingData.target}
              percentiles={props.whalingData.percentiles_open}
            />
          )}
          <RenderByAttribute
            attribute="Average Ships dropped by tier"
            data={props.whalingData.avg_by_attribute["tier"]}
          />
          <RenderByAttribute
            attribute="Average Rare/Not Rate Ships dropped"
            data={props.whalingData.avg_by_attribute["rare"]}
          />
        </Flex>
      </View>

      <Divider size="S" />
      <View>
        <Flex direction="row" gap="size-100" justifyContent="space-evenly" wrap>
          <RenderItems
            items={other_cat.resource}
            title="Resources"
            prefix="Average "
          />
          <RenderItems
            items={other_cat.eco}
            title="Economic Bonuses"
            prefix="Average "
          />
          <RenderItems
            items={other_cat.other}
            title="Other Items"
            prefix="Average "
          />
        </Flex>
      </View>
    </Flex>
  );
}

export default WhalingResults;
