"use client";

import { useEffect } from "react";
import { useTeamEntities } from "@/components/context/team-entities-context";

export function IncidentHydrator({
  incidentID,
  title,
}: {
  incidentID: string;
  title: string;
}) {
  const { setIncidentTitle } = useTeamEntities();

  useEffect(() => {
    setIncidentTitle(incidentID, title);
  }, [incidentID, title, setIncidentTitle]);

  return null;
}

