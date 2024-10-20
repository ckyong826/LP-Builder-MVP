import grapesjs, { Editor } from "grapesjs";
import GjsEditor, { Canvas } from "@grapesjs/react";

export default function CustomEditor() {
  const onEditor = (editor: Editor) => {
    console.log("Editor loaded", { editor });
  };

  return (
    <GjsEditor
      grapesjs={grapesjs}
      grapesjsCss="https://unpkg.com/grapesjs/dist/css/grapes.min.css"
      onEditor={onEditor}
      options={{
        height: "100vh",
        storageManager: false,
      }}
    >
      <div>
        // ... use your UI components
        <Canvas /> // place the GrapesJS canvas where you wish // ...
      </div>
    </GjsEditor>
  );
}
