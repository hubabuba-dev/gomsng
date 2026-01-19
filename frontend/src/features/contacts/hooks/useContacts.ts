import { useEffect, useState } from "react";
import { http } from "../../../shared/api/http";
import type { ContactSummary } from "../types";

export function useContacts() {
  const [data, setData] = useState<ContactSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function loadContacts() {
      try {
        const result = await http<ContactSummary[]>("/contacts", { auth: true });
        if (!cancelled) setData(Array.isArray(result) ? result : []);
      } catch {
        if (!cancelled) setData([]);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    loadContacts();

    return () => {
      cancelled = true;
    };
  }, []);

  return { data, loading };
}
