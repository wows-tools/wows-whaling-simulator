import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Image,
  Link,
  Grid,
  repeat,
  Divider,
  View,
  IllustratedMessage,
  Heading,
  Content,
  SearchField,
  Flex,
  Text,
  ActionGroup,
  Item,
  Badge,
} from "@adobe/react-spectrum";
import { Link as RouterLink } from "react-router-dom";

import { API_ROOT } from "../api-config";

const CATEGORY_ORDER = [
  "Santa's Gifts",
  "Black Friday",
  "Premium & Special Ships",
  "National Navies",
  "Collection",
  "Historical Events",
  "Collaborations",
  "Seasonal",
  "Anniversaries",
  "Leagues",
  "Recruitment",
  "Resources & Supplies",
  "Special Events",
];

// Distinct color per category so the badges are visually distinguishable.
const CATEGORY_COLOR = {
  "Santa's Gifts":           "seafoam",
  "Black Friday":            "fuchsia",
  "Premium & Special Ships": "indigo",
  "National Navies":         "blue",
  "Historical Events":       "orange",
  "Collection":              "celery",
  "Collaborations":          "purple",
  "Seasonal":                "green",
  "Anniversaries":           "magenta",
  "Leagues":                 "notice",
  "Recruitment":             "positive",
  "Resources & Supplies":    "gray",
  "Special Events":          "red",
};

export default function LootboxList() {
  const [lootboxes, setLootboxes] = useState([]);
  const [search, setSearch] = useState("");
  const [activeCategory, setActiveCategory] = useState(null);

  useEffect(() => {
    axios.get(`${API_ROOT}/api/v1/lootboxes`).then((res) => {
      setLootboxes(res.data["lootboxes"] || []);
    });
  }, []);

  // Build full grouped map (unfiltered) for the category buttons.
  const allGrouped = {};
  lootboxes.forEach((lb) => {
    const cat = lb.category || "Special Events";
    if (!allGrouped[cat]) allGrouped[cat] = [];
    allGrouped[cat].push(lb);
  });
  const allCategories = [
    ...CATEGORY_ORDER.filter((c) => allGrouped[c]),
    ...Object.keys(allGrouped).filter((c) => !CATEGORY_ORDER.includes(c)).sort(),
  ];

  // Apply search + active category filter.
  const needle = search.trim().toLowerCase();
  const filtered = lootboxes.filter((lb) => {
    const cat = lb.category || "Special Events";
    if (activeCategory && cat !== activeCategory) return false;
    if (needle && !lb.name.toLowerCase().includes(needle)) return false;
    return true;
  });

  const grouped = {};
  filtered.forEach((lb) => {
    const cat = lb.category || "Special Events";
    if (!grouped[cat]) grouped[cat] = [];
    grouped[cat].push(lb);
  });
  const visibleCategories = [
    ...CATEGORY_ORDER.filter((c) => grouped[c]),
    ...Object.keys(grouped).filter((c) => !CATEGORY_ORDER.includes(c)).sort(),
  ];

  return (
    <View margin="size-400">
      <Heading level={1}>World of Warships Lootbox Whaling Simulator</Heading>
      <Divider />
      <Heading level={2}>Container List</Heading>

      {/* Category filter buttons */}
      <ActionGroup
        selectionMode="single"
        selectedKeys={activeCategory ? [activeCategory] : []}
        onSelectionChange={(keys) => {
          const selected = [...keys][0] ?? null;
          setActiveCategory(selected);
        }}
        overflowMode="wrap"
        marginBottom="size-300"
      >
        {allCategories.map((cat) => (
          <Item key={cat}>
            <Badge
              variant="neutral"
              UNSAFE_style={{ marginRight: 6 }}
            >
              {allGrouped[cat].length}
            </Badge>
            <Text>{cat}</Text>
          </Item>
        ))}
      </ActionGroup>

      {/* Search */}
      <Flex marginBottom="size-400" maxWidth="size-3600">
        <SearchField
          label="Search containers"
          value={search}
          onChange={setSearch}
          onClear={() => setSearch("")}
          width="100%"
          aria-label="Search containers"
        />
      </Flex>

      {visibleCategories.length === 0 && (
        <Text>No containers match your search.</Text>
      )}

      {visibleCategories.map((cat) => (
        <View key={cat} marginBottom="size-600">
          <Flex alignItems="center" gap="size-100" marginBottom="size-200">
            <Badge variant="neutral">
              {grouped[cat].length}
            </Badge>
            <Heading level={3} margin="size-0">
              {cat}
            </Heading>
          </Flex>
          <Grid
            columns={repeat("auto-fit", "size-3600")}
            autoRows="size-3000"
            marginStart="size-400"
            gap="size-400"
          >
            {grouped[cat].map((lb) => (
              <View
                key={lb.id}
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
      ))}
    </View>
  );
}
