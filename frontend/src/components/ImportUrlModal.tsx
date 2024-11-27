import { useState } from "react";
import { Dialog } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { templateService } from "@/api/services/templateService";
import { BASE_URL } from "@/api/constants";

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
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const urlValue = formData.get("url") as string;

    if (!urlValue.trim()) {
      setError("Please enter a URL");
      return;
    }

    try {
      setLoading(true);
      setError("");

      const data = await templateService.convertUrl(urlValue);

      if (data.conversion.status === "failed") {
        throw new Error(
          data.conversion.error_message || "Failed to process URL"
        );
      }

      if (!data.html_path) {
        throw new Error("HTML path not found in the conversion result");
      }

      const htmlPath = data.html_path.replace("./output", "/output");
      const htmlContent = await templateService.fetchContent(htmlPath);

      if (!htmlContent) {
        throw new Error("Failed to fetch HTML content");
      }

      const filePaths = JSON.parse(data.file_paths);

      if (!filePaths || typeof filePaths !== "object") {
        throw new Error("Invalid file paths data");
      }

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
      (path) => `${BASE_URL}${path.replace("./output", "output")}`
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <div className="p-6">
        <h2 className="text-lg font-semibold mb-4">Import from URL</h2>
        <form onSubmit={handleSubmit}>
          <Input
            type="url"
            name="url"
            placeholder="Enter website URL"
            required
            className="mb-4"
          />
          {error && <p className="text-red-500 mb-4">{error}</p>}
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? <Spinner className="mr-2" /> : null}
              Import
            </Button>
          </div>
        </form>
      </div>
    </Dialog>
  );
}
