import { useState } from "react";
import { Dialog } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { templateService } from "@/api/services/templateService";
import { BASE_URL } from "@/api";

interface ImportUrlModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (content: {
    html: string;
    css: Record<string, string>;
    js: Record<string, string>;
    images: Record<string, string>;
  }) => void;
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

      const rawContent = await templateService.fetchContent(data.id);

      let content;
      try {
        content =
          typeof rawContent === "string" ? JSON.parse(rawContent) : rawContent;
      } catch (parseError) {
        console.error("JSON parse error:", parseError);
        throw new Error("Invalid response format from server");
      }

      if (!content.html) {
        throw new Error("Failed to fetch HTML content");
      }

      let processedHtml = content.html;
      const processedImages: Record<string, string> = {};

      if (content.images) {
        for (const [imagePath, imageUrl] of Object.entries(content.images)) {
          const url = imageUrl as string;
          try {
            // Clean up the URL: replace backslashes with forward slashes
            const cleanUrl = url.replace(/\\/g, "/");

            // Check if URL is already absolute
            const fullUrl = cleanUrl.startsWith("http")
              ? cleanUrl
              : `${BASE_URL}${cleanUrl.startsWith("/") ? "" : "/"}${cleanUrl}`;

            // Remove any duplicated URLs in the path
            const deduplicatedUrl = fullUrl.replace(
              /http:\/\/[^/]+(.*?)\1$/,
              "$1"
            );

            // Store the original path as key and the cleaned URL as value
            processedImages[imagePath] = deduplicatedUrl;
          } catch (error) {
            console.error("Failed to process image:", error);
          }
        }
      }
      // Combine all CSS files into one
      const combinedCss = Object.entries(content.css)
        .map(([filename, content]) => `/* ${filename} */\n${content}`)
        .join("\n\n");

      // Combine all JS files into one
      const combinedJs = Object.entries(content.js)
        .map(([filename, content]) => `// ${filename}\n${content}`)
        .join("\n\n");

      // Pass the processed content to the parent component
      onImport({
        html: processedHtml,
        css: { "combined.css": combinedCss },
        js: { "combined.js": combinedJs },
        images: processedImages,
      });

      onClose();
    } catch (err: any) {
      console.error("Import error:", err);
      setError(err.message || "An unexpected error occurred");
    } finally {
      setLoading(false);
    }
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
