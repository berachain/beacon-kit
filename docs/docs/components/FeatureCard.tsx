import React from "react";
import "../styles.css";

interface FeatureCardProps {
  title: string;
  description: string;
}

export default function FeatureCard({ title, description }: FeatureCardProps) {
  return (
    <div className="pl-2 pr-2 max-sm:px-0 max-lg:pb-3 max-lg:pr-0 w-1/4 max-lg:w-1/2 max-sm:w-full z-0">
      <div className="flex flex-col h-full w-full border border-border rounded-md p-6 max-lg:h-[142px]">
        {" "}
        <div className="text-md font-semibold">{title}</div>
        <div className="text-md text-muted-foreground">{description}</div>
      </div>
    </div>
  );
}
