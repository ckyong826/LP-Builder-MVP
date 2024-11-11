////////////////////////////////////////////
/////////// Import dependencies ////////////
////////////////////////////////////////////
"use client";

// GrapeJS
import grapesjs, { Editor } from "grapesjs";
import "grapesjs/dist/css/grapes.min.css";
import GjsEditor from "@grapesjs/react";
import gPreset from "grapesjs-preset-webpage";
import gBlocks from "grapesjs-blocks-basic";
import gCkEditor from "grapesjs-plugin-ckeditor";
import gFlexbox from "grapesjs-blocks-flexbox";
import gCustomCode from "grapesjs-custom-code";
import gNavbar from "grapesjs-navbar";
import gExport from "grapesjs-plugin-export";
import gStyleBg from "grapesjs-style-bg";
import gToolTip from "grapesjs-tooltip";

export default function CustomEditor() {
  const onEditor = (editor: Editor) => {
    console.log("Editor loaded", { editor });
  };

  return (
    <GjsEditor
      grapesjs={grapesjs}
      onEditor={onEditor}
      options={{
        height: "100vh",
        storageManager: false,
        container: "#gjs",
        fromElement: true,
      }}
      plugins={[
        gBlocks,
        gFlexbox,
        gPreset,
        gCustomCode,
        gNavbar,
        gExport,

        gStyleBg,
        gToolTip,
      ]}
    ></GjsEditor>
  );
}
