import { useEffect, useState } from "react";
import { http } from "../../../shared/api/http";
import type { ChatSummary } from "../types";

export function useChats() {
  const [data, setData] = useState<ChatSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function loadChats() {
      try {
        const result = await http<ChatSummary[]>("/chats", { auth: true });
        if (!cancelled) setData(Array.isArray(result) ? result : []);
      } catch {
        if (!cancelled) setData([]);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    loadChats();

    return () => {
      cancelled = true;
    };
  }, []);

  return { data, loading };
}
