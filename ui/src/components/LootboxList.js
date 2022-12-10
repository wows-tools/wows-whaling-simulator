import React from "react";
import axios from "axios";
import { Image } from "@adobe/react-spectrum";
import { Link } from "@adobe/react-spectrum";
import { Text } from "@adobe/react-spectrum";
import { Grid } from "@adobe/react-spectrum";
import { repeat } from "@adobe/react-spectrum";
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
        <h1>Pick a Containers:</h1>
        <Grid
          columns={repeat("auto-fit", "size-2400")}
          autoRows="size-3000"
          marginStart="size-800"
          gap="size-400"
        >
          {this.state.lootboxes.map((lb) => (
            <View
              backgroundColor="gray-200"
              borderColor="dark"
              borderRadius="small"
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
