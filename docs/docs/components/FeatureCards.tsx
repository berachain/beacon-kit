import React from "react";
import "../styles.css";
import FeatureCard from "./FeatureCard";
import Rails from "./Rails";
const featureCards = [
  {
    title: "Modular 🐻",
    description:
      "Composable modules to build applications & libraries with speed",
  },
  {
    title: "Lightweight 🦋",
    description: "Tiny bundle size optimized for tree-shaking",
  },
  {
    title: "Performant 🦄",
    description: "Optimized architecture compared to alternative libraries",
  },
  {
    title: "Customizable 🦥",
    description:
      "Designed from the ground up to be easily customized to your needs",
  },
];
export default function FeatureCards() {
  return (
    <div className="flex justify-between flex-wrap">
      {featureCards.map((card) => (
        <FeatureCard title={card.title} description={card.description} />
      ))}
    </div>
  );
}
