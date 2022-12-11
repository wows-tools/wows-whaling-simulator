import React from "react";
import axios from "axios";
import { Image } from "@adobe/react-spectrum";
import { Link } from "@adobe/react-spectrum";
import { Text } from "@adobe/react-spectrum";
import { Grid } from "@adobe/react-spectrum";
import { repeat } from "@adobe/react-spectrum";
import { Flex } from "@adobe/react-spectrum";
import { Divider } from "@adobe/react-spectrum";
import { View } from "@adobe/react-spectrum";
import { IllustratedMessage } from "@adobe/react-spectrum";
import { Heading } from "@adobe/react-spectrum";
import { Content } from "@adobe/react-spectrum";
import { Link as RouterLink } from "react-router-dom";

import { API_ROOT } from "../api-config";

export default class LootboxList extends React.Component {
  state = {
    lootboxes: [],
  };

  componentDidMount() {
    axios.get(`${API_ROOT}/api/v1/lootboxes`).then((res) => {
      const lootboxes = res.data["lootboxes"];
      this.setState({ lootboxes });
    });
  }

  render() {
    return (
      <View marging="size-400">
        <Heading level={1}>World of Warships Lootbox Whaling Simulator</Heading>
        <Divider />
        <Heading level={2}>Disclaimer</Heading>
        <View
          backgroundColor="red-400"
          colorVersion="6"
          width="size-6000"
          padding="size-100"
          borderColor="red-800"
          borderWidth="medium"
          borderRadius="medium"
        >
          <Flex gap="size-100" direction="column">
            <Text>
              This website is in alpha stage. You are likely to encounter bugs.
            </Text>
            <Text>
              The estimates are provided as is, without an warranty of any kind.
            </Text>
          </Flex>
        </View>

        <Heading level={2}>Container List</Heading>

        <Heading level={3}>Pick a Containers:</Heading>
        <Grid
          columns={repeat("auto-fit", "size-3600")}
          autoRows="size-3000"
          marginStart="size-800"
          gap="size-400"
        >
          {this.state.lootboxes.map((lb) => (
            <View
              width="size-3600"
              backgroundColor="gray-100"
              borderRadius="medium"
              borderWidth="thin"
              borderColor="dark"
              padding="size-100"
            >
              <Link>
                <RouterLink to={"/lootboxes/" + lb.id}>
                  <IllustratedMessage>
                    <Image
                      height="200px"
                      objectFit="scale-down"
                      src={API_ROOT + lb.img}
                      alt={lb.name}
                    />
                    <Content>{lb.name}</Content>
                  </IllustratedMessage>
                </RouterLink>
              </Link>
            </View>
          ))}
        </Grid>
      </View>
    );
  }
}
