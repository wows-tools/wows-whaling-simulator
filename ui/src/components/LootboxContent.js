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

function checkUnset(props) {
  return props === undefined || props === null || props.length === 0;
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
      <Flex direction="row" gap="size-100" justifyContent="space-evenly" wrap>
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
            footer={
              "(ie: €" +
              Number(lootbox.price / lootbox.exchange_rate_euro).toFixed(1) +
              " or $" +
              Number(lootbox.price / lootbox.exchange_rate_dollar).toFixed(2) +
              ")"
            }
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
      </Flex>
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

export default LootboxContent;
