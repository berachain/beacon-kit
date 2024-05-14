import React from "react";
import "../styles.css";

export default function Overview() {
  return (
    <>
      <div className="voc_Div">
        <header className="vocs_Header border-dotted">
          <h1 className="vocs_H1 vocs_Heading">Features</h1>
        </header>
        <div className="flex flex-col w-full mx-auto md:flex-row gap-8">
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md">
            <header className="vocs_Header border-none flex items-center flex-col">
              <svg
                width="100"
                height="100"
                viewBox="0 0 398 400"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <rect
                  x="2"
                  y="2"
                  width="196"
                  height="396"
                  stroke="#D57F00"
                  stroke-width="4"
                  stroke-dasharray="12.8 12.8"
                />
                <rect
                  x="200"
                  y="2"
                  width="196"
                  height="396"
                  stroke="#D57F00"
                  stroke-width="4"
                  stroke-dasharray="12.8 12.8"
                />
              </svg>

              <h1 className="vocs_H2 vocs_Heading">What is BeaconKit</h1>
            </header>

            <p className="vocs_Paragraph">
              BeaconKit introduces an innovative framework that utilizes the
              Cosmos-SDK to create a flexible, customizable consensus layer
              tailored for Ethereum-based blockchains
            </p>
          </div>
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md">
            <header className="vocs_Header border-none flex items-center flex-col">
              <svg
                width="100"
                height="100"
                viewBox="0 0 398 400"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <rect
                  x="2"
                  y="2"
                  width="196"
                  height="396"
                  stroke="#D57F00"
                  stroke-width="4"
                  stroke-dasharray="12.8 12.8"
                />
                <rect
                  x="200"
                  y="2"
                  width="196"
                  height="396"
                  stroke="#D57F00"
                  stroke-width="4"
                  stroke-dasharray="12.8 12.8"
                />
                <path
                  d="M379.25 0H391.75V2H393.172L392.04 3.1312L393.455 4.54541L384.364 13.6363L382.949 12.2221L373.859 21.313L375.273 22.7272L366.182 31.8181L364.768 30.4039L355.677 39.4948L357.091 40.909L348 50L346.586 48.5858L337.495 57.6767L338.909 59.0909L329.818 68.1818L328.404 66.7676L319.313 75.8585L320.727 77.2727L311.636 86.3636L310.222 84.9494L301.131 94.0403L302.545 95.4545L293.455 104.545L292.04 103.131L282.949 112.222L284.364 113.636L275.273 122.727L273.859 121.313L264.768 130.404L266.182 131.818L257.091 140.909L255.677 139.495L246.586 148.586L248 150L238.909 159.091L237.495 157.677L228.404 166.768L229.818 168.182L220.727 177.273L219.313 175.859L210.222 184.949L211.636 186.364L202.545 195.455L201.131 194.04L200 195.172V193.75H198V181.25H200V168.75H198V156.25H200V143.75H198V131.25H200V118.75H198V106.25H200V93.75H198V81.25H200V68.75H198V56.25H200V43.75H198V31.25H200V18.75H198V6.25H200V2H204.25V0H216.75V2H229.25V0H241.75V2H254.25V0H266.75V2H279.25V0H291.75V2H304.25V0H316.75V2H329.25V0H341.75V2H354.25V0H366.75V2H379.25V0ZM198 218.75V206.25H200V204.828L201.131 205.96L202.545 204.545L211.636 213.636L210.222 215.051L219.313 224.141L220.727 222.727L229.818 231.818L228.404 233.232L237.495 242.323L238.909 240.909L248 250L246.586 251.414L255.677 260.505L257.091 259.091L266.182 268.182L264.768 269.596L273.858 278.687L275.273 277.273L284.364 286.364L282.949 287.778L292.04 296.869L293.455 295.455L302.545 304.545L301.131 305.96L310.222 315.051L311.636 313.636L320.727 322.727L319.313 324.141L328.404 333.232L329.818 331.818L338.909 340.909L337.495 342.323L346.586 351.414L348 350L357.091 359.091L355.677 360.505L364.768 369.596L366.182 368.182L375.273 377.273L373.859 378.687L382.949 387.778L384.364 386.364L393.455 395.455L392.04 396.869L393.172 398H391.75V400H379.25V398H366.75V400H354.25V398H341.75V400H329.25V398H316.75V400H304.25V398H291.75V400H279.25V398H266.75V400H254.25V398H241.75V400H229.25V398H216.75V400H204.25V398H200V393.75H198V381.25H200V368.75H198V356.25H200V343.75H198V331.25H200V318.75H198V306.25H200V293.75H198V281.25H200V268.75H198V256.25H200V243.75H198V231.25H200V218.75H198Z"
                  stroke="#D57F00"
                  stroke-width="4"
                  stroke-dasharray="12.8 12.8"
                />
                <g opacity="0.5">
                  <mask id="path-4-inside-1_218_11" fill="white">
                    <path d="M0 0H398V200H0V0Z" />
                  </mask>
                  <path
                    d="M0 202H7.10714V198H0V202ZM21.3214 202H35.5357V198H21.3214V202ZM49.75 202H63.9643V198H49.75V202ZM78.1786 202H92.3929V198H78.1786V202ZM106.607 202H120.821V198H106.607V202ZM135.036 202H149.25V198H135.036V202ZM163.464 202H177.679V198H163.464V202ZM191.893 202H206.107V198H191.893V202ZM220.321 202H234.536V198H220.321V202ZM248.75 202H262.964V198H248.75V202ZM277.179 202H291.393V198H277.179V202ZM305.607 202H319.821V198H305.607V202ZM334.036 202H348.25V198H334.036V202ZM362.464 202H376.679V198H362.464V202ZM390.893 202H398V198H390.893V202ZM0 204H7.10714V196H0V204ZM21.3214 204H35.5357V196H21.3214V204ZM49.75 204H63.9643V196H49.75V204ZM78.1786 204H92.3929V196H78.1786V204ZM106.607 204H120.821V196H106.607V204ZM135.036 204H149.25V196H135.036V204ZM163.464 204H177.679V196H163.464V204ZM191.893 204H206.107V196H191.893V204ZM220.321 204H234.536V196H220.321V204ZM248.75 204H262.964V196H248.75V204ZM277.179 204H291.393V196H277.179V204ZM305.607 204H319.821V196H305.607V204ZM334.036 204H348.25V196H334.036V204ZM362.464 204H376.679V196H362.464V204ZM390.893 204H398V196H390.893V204Z"
                    fill="#D57F00"
                    mask="url(#path-4-inside-1_218_11)"
                  />
                </g>
              </svg>

              <h1 className="vocs_H2 vocs_Heading">How BeaconKit works</h1>
            </header>
            <p className="vocs_Paragraph">
              BeaconKit introduces an innovative framework that utilizes the
              Cosmos-SDK to create a flexible, customizable consensus layer
              tailored for Ethereum-based blockchains
            </p>
          </div>
          <div className="flex flex-col items-center justify-between border-border border p-5 rounded-md">
            <header className="vocs_Header border-none flex items-center flex-col">
              <svg
                width="100"
                height="100"
                viewBox="0 0 398 400"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <rect width="200" height="400" fill="#D57F00" />
                <path d="M198 0H398L198 200L398 400H198V0Z" fill="#D57F00" />
              </svg>

              <h1 className="vocs_H2 vocs_Heading">Why BeaconKit</h1>
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
