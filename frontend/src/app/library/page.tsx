"use client";

import { useState } from "react";
import grapesjs, { Editor } from "grapesjs";
import "grapesjs/dist/css/grapes.min.css";
import GjsEditor from "@grapesjs/react";
import gPreset from "grapesjs-preset-webpage";
import gBlocks from "grapesjs-blocks-basic";
import gFlexbox from "grapesjs-blocks-flexbox";
import gCustomCode from "grapesjs-custom-code";
import gNavbar from "grapesjs-navbar";
import gExport from "grapesjs-plugin-export";
import gStyleBg from "grapesjs-style-bg";
import gToolTip from "grapesjs-tooltip";
import { ImportUrlModal } from "@/components/ImportUrlModal";

export default function CustomEditor() {
  const [editor, setEditor] = useState<Editor | null>(null);
  const [isImportModalOpen, setIsImportModalOpen] = useState(false);

  const onEditor = (editor: Editor) => {
    setEditor(editor);

    // Add import button to panel
    editor.Panels.addButton("options", {
      id: "import-url",
      className: "fa fa-download",
      command: "open-import-modal",
      attributes: { title: "Import from URL" },
    });

    // Add command for opening import modal
    editor.Commands.add("open-import-modal", {
      run: () => setIsImportModalOpen(true),
      stop: () => setIsImportModalOpen(false),
    });
  };

  const handleImport = async (content: {
    html: string;
    css: Record<string, string>;
    js: Record<string, string>;
    images: Record<string, Uint8Array<ArrayBufferLike>>;
  }) => {
    if (!editor) return;

    try {
      // Clear existing content
      editor.setComponents("");
      editor.setStyle("");

      // Load CSS files
      const cssPromises = Object.entries(content.css).map(
        async ([cssPath, cssContent]) => {
          editor.setStyle(cssContent);
        }
      );

      // Load JS files
      const jsPromises = Object.entries(content.js).map(
        async ([jsPath, jsContent]) => {
          editor.Components.addComponent({
            type: "script",
            content: jsContent,
          });
        }
      );

      // Process images in HTML
      Object.entries(content.images).forEach(([imagePath, imageContent]) => {
        const filename = imagePath.split("/").pop();
        // Replace image paths in HTML content
        content.html = content.html.replace(
          new RegExp(filename as string, "g"),
          imagePath
        );
      });

      // Load HTML content
      editor.setComponents(content.html);

      // Wait for all assets to load
      await Promise.all([...cssPromises, ...jsPromises]);

      // Refresh the editor
      editor.refresh();
    } catch (error) {
      console.error("Error importing content:", error);
      // You might want to add error handling UI here
    }
  };

  return (
    <>
      <GjsEditor
        grapesjs={grapesjs}
        onEditor={onEditor}
        options={{
          height: "100vh",
          storageManager: false,
          container: "#gjs",
          fromElement: true,
          assetManager: {
            upload: false,
            assets: [],
          },
          styleManager: {
            sectors: [
              /* your style sectors */
            ],
          },
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
      />

      <ImportUrlModal
        isOpen={isImportModalOpen}
        onClose={() => setIsImportModalOpen(false)}
        onImport={handleImport}
      />
    </>
  );
}
