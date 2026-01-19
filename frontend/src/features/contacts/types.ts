export type ContactSummary = {
  id: string;
  name: string;
  status: "online" | "away" | "offline";
  lastSeen?: string;
};
