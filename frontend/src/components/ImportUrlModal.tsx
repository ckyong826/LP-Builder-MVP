import { useState } from "react";
import { Dialog } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { templateService } from "@/api/services/templateService";
import { BASE_URL } from "@/api/constants";
import type { ConvertUrlResponse } from "@/types/models";

interface ImportUrlModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (
    htmlContent: string,
    assets: { css: string[]; js: string[]; images: string[] }
  ) => void;
}

export function ImportUrlModal({
  isOpen,
  onClose,
  onImport,
}: ImportUrlModalProps) {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleImport = async () => {
    try {
      setLoading(true);
      setError("");

      // Use the template service to convert URL
      const data = await templateService.convertUrl(url);

      if (data.conversion.status === "failed") {
        throw new Error(
          data.conversion.error_message || "Failed to process URL"
        );
      }

      // Fetch HTML content
      const htmlPath = data.html_path.replace("./output", "/output");
      const htmlContent = await templateService.fetchContent(htmlPath);

      // Process file paths
      const filePaths = JSON.parse(data.file_paths);
      const convertedFilePaths = {
        css: processFilePaths(filePaths.css),
        js: processFilePaths(filePaths.js),
        images: processFilePaths(filePaths.images),
      };

      onImport(htmlContent, convertedFilePaths);
      onClose();
    } catch (err: any) {
      console.error("Import error:", err);
      setError(err.message || "An unexpected error occurred");
    } finally {
      setLoading(false);
    }
  };

  // Helper function to process file paths
  const processFilePaths = (paths: string[] = []): string[] => {
    return paths.map(
      (path) => `${BASE_URL}${path.replace("./output", "/output")}`
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <div className="p-6">
        <h2 className="text-lg font-semibold mb-4">Import from URL</h2>
        <Input
          type="url"
          placeholder="Enter website URL"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          className="mb-4"
        />
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <div className="flex justify-end gap-2">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleImport} disabled={loading || !url.trim()}>
            {loading ? <Spinner className="mr-2" /> : null}
            Import
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
