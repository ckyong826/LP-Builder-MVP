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
  onImport: (content: {
    html: string;
    css: Record<string, string>;
    js: Record<string, string>;
    images: Record<string, Uint8Array>;
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

      const content = JSON.parse(await templateService.fetchContent(data.id));

      if (!content.html) {
        throw new Error("Failed to fetch HTML content");
      }

      onImport({
        html: content.html,
        css: content.css,
        js: content.js,
        images: Object.fromEntries(
          Object.entries(content.images).map(([key, value]) => [
            key,
            new Uint8Array(Object.values(value as number[])),
          ])
        ),
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
