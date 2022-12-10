import React, { Component } from "react";
import styled from "styled-components";
import LinesEllipsis from "react-lines-ellipsis";
import { ClipLoader } from "react-spinners";
import { View } from "@adobe/react-spectrum";
import { Text } from "@adobe/react-spectrum";
import { Heading } from "@adobe/react-spectrum";

import NumericContent from "./Numeric";

const SubHeader = styled.div`
  padding: 2px 0;
  font-size: 12px;
  color: #7a7c7d;
  height: 2.4em;
`;

const Loading = styled.div`
  position: absolute;
  width: 9rem;
  height: 9rem;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: -1;
`;

const TileContent = styled.div`
  height: 5rem;
`;

const Footer = styled.div`
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 13px;
`;

const linesEllipsisSettings = {
  ellipsis: "...",
  trimRight: true,
  basedOn: "letters",
};

export default class ReactGenericTile extends Component {
  render() {
    const {
      header,
      subheader,
      footer,
      number,
      scale,
      indicator,
      loading,
      icon,
      size,
      onClick,
      minHeight,
      minWidth,
    } = this.props;
    let { color } = this.props;

    switch (color) {
      case "Good":
        color = "#2b7c2b";
        break;
      case "Warning":
        color = "#e78c07";
        break;
      case "Bad":
        color = "#bb0000";
        break;
      default:
        break;
    }

    return (
      <View
        minHeight={minHeight}
        minWidth={minWidth}
        backgroundColor="gray-100"
        borderRadius="medium"
        borderWidth="thin"
        borderColor="dark"
        padding="size-100"
      >
        <Text>
          <ClipLoader
            sizeUnit={"px"}
            size={44}
            color={"#123abc"}
            loading={loading}
          />
        </Text>

        <Heading level={2}>
          <LinesEllipsis text={header} maxLine="2" {...linesEllipsisSettings} />
        </Heading>
        <Heading
          level={5}
          UNSAFE_style={{
            "font-weight": "normal",
            color: "gray",
            "font-style": "italic",
          }}
        >
          <LinesEllipsis
            text={subheader}
            maxLine="2"
            {...linesEllipsisSettings}
          />
        </Heading>

        <TileContent>
          <NumericContent
            icon={icon}
            number={number}
            scale={scale}
            color={color}
            indicator={indicator}
          />
        </TileContent>

        <Footer>{footer}</Footer>
      </View>
    );
  }
}

// <TileContent footer="Current Quarter" unit="EUR" class="sapUiSmallMargin">
//   <content>
//     <NumericContent scale="M" value="1.96"
//       valueColor="Error" indicator="Up" />
//   </content>
// </TileContent>

ReactGenericTile.defaultProps = {
  minHeight: "size-2000",
  minWidth: "size-2000",
  header: "",
  subheader: "",
  footer: "",
  loading: false,
  number: "",
  scale: "",
  color: "#000",
  indicator: null,
  icon: null,
  size: "Normal",
};
