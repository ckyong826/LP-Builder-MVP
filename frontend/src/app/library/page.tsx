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
    images: Record<string, string>; // Changed type to string URLs
  }) => {
    if (!editor) return;

    try {
      // Clear existing content
      editor.setComponents("");
      editor.setStyle("");

      // Add images to asset manager and keep track of their new URLs
      const assetManager = editor.AssetManager;
      const imageUrlMap = new Map<string, string>();

      // Process all images
      for (const [originalUrl, imageUrl] of Object.entries(content.images)) {
        // Add image to asset manager
        assetManager.add({
          src: imageUrl,
          type: "image",
        });

        // Store the mapping of original URL to new URL
        imageUrlMap.set(originalUrl, imageUrl);
      }

      // Process HTML to replace image URLs
      let processedHtml = content.html;
      imageUrlMap.forEach((newUrl, originalUrl) => {
        // Get filename from URL by removing path and query parameters
        const filename = originalUrl.split("/").pop()?.split("?")[0];
        if (!filename) return;

        // Create a pattern that matches src attribute with any content
        const pattern = new RegExp(`src=["'][^"']*${filename}[^"']*["']`, "g");

        // Replace all occurrences
        processedHtml = processedHtml.replace(pattern, `src="${newUrl}"`);
      });
      // Add combined CSS
      const cssContent = Object.values(content.css)[0];
      editor.setStyle(cssContent);

      // Add combined JS
      const jsContent = Object.values(content.js)[0];
      editor.Components.addComponent({
        type: "script",
        content: jsContent,
      });

      // Load processed HTML content
      editor.setComponents(processedHtml);

      // Refresh the editor
      editor.refresh();
    } catch (error) {
      console.error("Error importing content:", error);
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
          deviceManager: {
            devices: [
              {
                name: "Desktop",
                width: "", // Default width
              },
              {
                name: "Mobile",
                width: "400px", // Adjust this value as needed
                widthMedia: "375px", // This sets the CSS media query width
              },
            ],
          },
          assetManager: {
            upload: false,
            assets: [],
            embedAsBase64: true, // Add this line to enable base64 support
            dropzone: true, // Optional: enables drag and drop of base64 images
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
