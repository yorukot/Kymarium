"use client";

import * as React from "react";
import { apiRequest } from "@/lib/api/client";

type MonitorSummary = {
  id: string;
  name: string;
};

type TeamEntitiesContextValue = {
  getMonitorName: (monitorID: string) => string | undefined;
  setMonitorName: (monitorID: string, name: string) => void;
  setMonitors: (monitors: MonitorSummary[]) => void;
  ensureMonitorLoaded: (monitorID: string) => void;
};

const TeamEntitiesContext = React.createContext<TeamEntitiesContextValue | null>(
  null,
);

export function TeamEntitiesProvider({
  teamID,
  children,
}: {
  teamID: string;
  children: React.ReactNode;
}) {
  const [monitorsByID, setMonitorsByID] = React.useState<
    Record<string, MonitorSummary>
  >({});

  const monitorsRef = React.useRef(monitorsByID);
  React.useEffect(() => {
    monitorsRef.current = monitorsByID;
  }, [monitorsByID]);

  const inflightRef = React.useRef<Set<string>>(new Set());

  const setMonitorName = React.useCallback((monitorID: string, name: string) => {
    const trimmed = name.trim();
    if (!trimmed) return;

    setMonitorsByID((prev) => {
      const existing = prev[monitorID];
      if (existing?.name === trimmed) return prev;
      return {
        ...prev,
        [monitorID]: {
          id: monitorID,
          name: trimmed,
        },
      };
    });
  }, []);

  const setMonitors = React.useCallback((monitors: MonitorSummary[]) => {
    setMonitorsByID((prev) => {
      let next: Record<string, MonitorSummary> | null = null;

      for (const monitor of monitors) {
        if (!monitor?.id) continue;
        const name = typeof monitor.name === "string" ? monitor.name.trim() : "";
        if (!name) continue;

        const existing = prev[monitor.id];
        if (existing?.name === name) continue;

        if (!next) next = { ...prev };
        next[monitor.id] = { id: monitor.id, name };
      }

      return next ?? prev;
    });
  }, []);

  const getMonitorName = React.useCallback(
    (monitorID: string) => monitorsByID[monitorID]?.name,
    [monitorsByID],
  );

  const ensureMonitorLoaded = React.useCallback(
    (monitorID: string) => {
      if (!monitorID) return;
      if (monitorsRef.current[monitorID]?.name) return;
      if (inflightRef.current.has(monitorID)) return;

      inflightRef.current.add(monitorID);

      void apiRequest<{ data?: { monitor?: { id?: string; name?: string } } }>(
        `/api/teams/${teamID}/monitors/${monitorID}/analytics`,
        {
          defaultError: "Failed to load monitor",
          redirectOn401: true,
        },
      )
        .then((res) => {
          const monitor = res.data?.data?.monitor;
          if (monitor?.id && typeof monitor.name === "string") {
            setMonitorName(monitor.id, monitor.name);
          }
        })
        .finally(() => {
          inflightRef.current.delete(monitorID);
        });
    },
    [setMonitorName, teamID],
  );

  const value = React.useMemo<TeamEntitiesContextValue>(
    () => ({
      getMonitorName,
      setMonitorName,
      setMonitors,
      ensureMonitorLoaded,
    }),
    [ensureMonitorLoaded, getMonitorName, setMonitorName, setMonitors],
  );

  return (
    <TeamEntitiesContext.Provider value={value}>
      {children}
    </TeamEntitiesContext.Provider>
  );
}

export function useTeamEntities() {
  const context = React.useContext(TeamEntitiesContext);
  if (!context) {
    throw new Error("useTeamEntities must be used within a TeamEntitiesProvider");
  }
  return context;
}
