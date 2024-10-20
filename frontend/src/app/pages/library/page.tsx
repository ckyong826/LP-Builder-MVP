////////////////////////////////////////////
/////////// Import dependencies ////////////
////////////////////////////////////////////
"use client";

// GrapeJS
import grapesjs, { Editor } from "grapesjs";
import "grapesjs/dist/css/grapes.min.css";
import GjsEditor, { Canvas } from "@grapesjs/react";
import gPreset from "grapesjs-preset-webpage";
import gBlocks from "grapesjs-blocks-basic";
import gFlexbox from "grapesjs-blocks-flexbox";
import gCountDown from "grapesjs-component-countdown";
import gCustomCode from "grapesjs-custom-code";
import gNavbar from "grapesjs-navbar";
import gCkEditor from "grapesjs-plugin-ckeditor";
import gExport from "grapesjs-plugin-export";
import gForms from "grapesjs-plugin-forms";
import gStyleBg from "grapesjs-style-bg";
import gStyleGradient from "grapesjs-style-gradient";
import gTooltip from "grapesjs-tooltip";

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
        gCountDown,
        gCustomCode,
        gNavbar,
        gCkEditor,
        gExport,
        gForms,
        gStyleBg,
        gStyleGradient,
        gTooltip,
      ]}
    >
      <div>
        <Canvas />
      </div>
    </GjsEditor>
  );
}
