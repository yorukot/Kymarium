"use client";

import * as React from "react";
import { apiRequest } from "@/lib/api/client";

type MonitorSummary = {
  id: string;
  name: string;
};

type IncidentSummary = {
  id: string;
  title: string;
};

type TeamEntitiesContextValue = {
  getMonitorName: (monitorID: string) => string | undefined;
  setMonitorName: (monitorID: string, name: string) => void;
  setMonitors: (monitors: MonitorSummary[]) => void;
  ensureMonitorLoaded: (monitorID: string) => void;

  getIncidentTitle: (incidentID: string) => string | undefined;
  setIncidentTitle: (incidentID: string, title: string) => void;
  setIncidents: (incidents: IncidentSummary[]) => void;
  ensureIncidentLoaded: (incidentID: string) => void;
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

  const [incidentsByID, setIncidentsByID] = React.useState<
    Record<string, IncidentSummary>
  >({});

  const incidentsRef = React.useRef(incidentsByID);
  React.useEffect(() => {
    incidentsRef.current = incidentsByID;
  }, [incidentsByID]);

  const inflightIncidentsRef = React.useRef<Set<string>>(new Set());

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

  const setIncidentTitle = React.useCallback(
    (incidentID: string, title: string) => {
      const trimmed = title.trim();
      if (!trimmed) return;

      setIncidentsByID((prev) => {
        const existing = prev[incidentID];
        if (existing?.title === trimmed) return prev;
        return {
          ...prev,
          [incidentID]: {
            id: incidentID,
            title: trimmed,
          },
        };
      });
    },
    [],
  );

  const setIncidents = React.useCallback((incidents: IncidentSummary[]) => {
    setIncidentsByID((prev) => {
      let next: Record<string, IncidentSummary> | null = null;

      for (const incident of incidents) {
        if (!incident?.id) continue;
        const title =
          typeof incident.title === "string" ? incident.title.trim() : "";
        if (!title) continue;

        const existing = prev[incident.id];
        if (existing?.title === title) continue;

        if (!next) next = { ...prev };
        next[incident.id] = { id: incident.id, title };
      }

      return next ?? prev;
    });
  }, []);

  const getIncidentTitle = React.useCallback(
    (incidentID: string) => incidentsByID[incidentID]?.title,
    [incidentsByID],
  );

  const ensureIncidentLoaded = React.useCallback(
    (incidentID: string) => {
      if (!incidentID) return;
      if (incidentsRef.current[incidentID]?.title) return;
      if (inflightIncidentsRef.current.has(incidentID)) return;

      inflightIncidentsRef.current.add(incidentID);

      void apiRequest<{ data?: { id?: string; title?: string | null } }>(
        `/api/teams/${teamID}/incidents/${incidentID}`,
        {
          defaultError: "Failed to load incident",
          redirectOn401: true,
        },
      )
        .then((res) => {
          const incident = res.data?.data;
          if (incident?.id && typeof incident.title === "string") {
            setIncidentTitle(incident.id, incident.title);
          }
        })
        .finally(() => {
          inflightIncidentsRef.current.delete(incidentID);
        });
    },
    [setIncidentTitle, teamID],
  );

  const value = React.useMemo<TeamEntitiesContextValue>(
    () => ({
      getMonitorName,
      setMonitorName,
      setMonitors,
      ensureMonitorLoaded,

      getIncidentTitle,
      setIncidentTitle,
      setIncidents,
      ensureIncidentLoaded,
    }),
    [
      ensureIncidentLoaded,
      ensureMonitorLoaded,
      getIncidentTitle,
      getMonitorName,
      setIncidentTitle,
      setIncidents,
      setMonitorName,
      setMonitors,
    ],
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
