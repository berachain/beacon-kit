import React from "react";
import "../styles.css";

export default function Overview() {
  return (
    <>
      <div className="voc_Div">
        <header className="vocs_Header border-dotted">
          <h1 className="vocs_H1 vocs_Heading">Features</h1>
        </header>
        <div className="flex flex-col max-w-3xl mx-auto md:flex-row gap-8">
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md shadow-md">
            <header className="vocs_Header border-none">
              <h1 className="vocs_H1 vocs_Heading">What is BeaconKit</h1>
            </header>

            <p className="vocs_Paragraph">
              BeaconKit introduces an innovative framework that utilizes the
              Cosmos-SDK to create a flexible, customizable consensus layer
              tailored for Ethereum-based blockchains
            </p>
          </div>
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md shadow-md">
            <header className="vocs_Header border-none">
              <h1 className="vocs_H1 vocs_Heading">How BeaconKit works</h1>
            </header>
            <p className="vocs_Paragraph">
              BeaconKit introduces an innovative framework that utilizes the
              Cosmos-SDK to create a flexible, customizable consensus layer
              tailored for Ethereum-based blockchains
            </p>
          </div>
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md shadow-md">
            <header className="vocs_Header border-none">
              <h1 className="vocs_H1 vocs_Heading">Why BeaconKit</h1>
            </header>
            <p className="vocs_Paragraph">
              BeaconKit introduces an innovative framework that utilizes the
              Cosmos-SDK to create a flexible, customizable consensus layer
              tailored for Ethereum-based blockchains
            </p>
          </div>
        </div>
        <header className="vocs_Header border-dotted">
          <h1 className="vocs_H1 vocs_Heading">Sponsors</h1>
        </header>
      </div>
    </>
  );
}
