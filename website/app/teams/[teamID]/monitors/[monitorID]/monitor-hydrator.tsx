"use client";

import { useEffect } from "react";
import { useTeamEntities } from "@/components/context/team-entities-context";

export function MonitorHydrator({
  monitorID,
  name,
}: {
  monitorID: string;
  name: string;
}) {
  const { setMonitorName } = useTeamEntities();

  useEffect(() => {
    setMonitorName(monitorID, name);
  }, [monitorID, name, setMonitorName]);

  return null;
}
