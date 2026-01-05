"use client";

import * as React from "react";

export type TeamSummary = {
  id: string;
  name: string;
  role: string;
};

const TeamsContext = React.createContext<TeamSummary[] | null>(null);

export function TeamsProvider({
  teams,
  children,
}: {
  teams: TeamSummary[];
  children: React.ReactNode;
}) {
  return (
    <TeamsContext.Provider value={teams}>{children}</TeamsContext.Provider>
  );
}

export function useTeams() {
  const context = React.useContext(TeamsContext);
  if (!context) {
    throw new Error("useTeams must be used within a TeamsProvider");
  }
  return context;
}
