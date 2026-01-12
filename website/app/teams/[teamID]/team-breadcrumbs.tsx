"use client";

import { Fragment, useEffect } from "react";
import { useParams, usePathname } from "next/navigation";
import {
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useTeamEntities } from "@/components/context/team-entities-context";
import { humanizeIdentifier } from "@/lib/parsers/strings";

const LABELS: Record<string, string> = {
  monitors: "Monitors",
  incidents: "Incidents",
  notifications: "Notifications",
  new: "New",
  view: "View",
  edit: "Edit",
};

export function TeamBreadcrumbs() {
  const { teamID } = useParams<{ teamID: string }>();
  const pathname = usePathname();
  const {
    ensureMonitorLoaded,
    getMonitorName,
    ensureIncidentLoaded,
    getIncidentTitle,
  } = useTeamEntities();

  const rawSegments = pathname.split("/").filter(Boolean);
  const segments = rawSegments[0] === "teams" ? rawSegments.slice(2) : rawSegments;

  const [section, idOrAction, action] = segments.length
    ? segments
    : (["monitors"] as const);

  const activeMonitorID =
    section === "monitors" && idOrAction && idOrAction !== "new"
      ? idOrAction
      : null;
  const activeMonitorName = activeMonitorID
    ? getMonitorName(activeMonitorID)
    : undefined;

  const activeIncidentID =
    section === "incidents" && idOrAction && idOrAction !== "new"
      ? idOrAction
      : null;
  const activeIncidentTitle = activeIncidentID
    ? getIncidentTitle(activeIncidentID)
    : undefined;

  useEffect(() => {
    if (!activeMonitorID) return;
    if (activeMonitorName) return;
    ensureMonitorLoaded(activeMonitorID);
  }, [activeMonitorID, activeMonitorName, ensureMonitorLoaded]);

  useEffect(() => {
    if (!activeIncidentID) return;
    if (activeIncidentTitle) return;
    ensureIncidentLoaded(activeIncidentID);
  }, [activeIncidentID, activeIncidentTitle, ensureIncidentLoaded]);

  const crumbs: Array<{ label: string; href?: string }> = [];

  const sectionLabel = LABELS[section] ?? humanizeIdentifier(section);
  crumbs.push({
    label: sectionLabel,
    href: idOrAction ? `/teams/${teamID}/${section}` : undefined,
  });

  if (section === "monitors") {
    if (idOrAction === "new") {
      crumbs.push({ label: LABELS.new });
    } else if (idOrAction) {
      const monitorID = idOrAction;
      crumbs.push({
        label: activeMonitorName ?? "Monitor",
        href: action ? `/teams/${teamID}/monitors/${monitorID}` : undefined,
      });

      if (action) {
        crumbs.push({ label: LABELS[action] ?? humanizeIdentifier(action) });
      }
    }
  } else if (section === "incidents") {
    if (idOrAction === "new") {
      crumbs.push({ label: LABELS.new });
    } else if (idOrAction) {
      const incidentID = idOrAction;
      crumbs.push({
        label: activeIncidentTitle ?? "Incident",
        href: action ? `/teams/${teamID}/incidents/${incidentID}` : undefined,
      });

      if (action) {
        crumbs.push({ label: LABELS[action] ?? humanizeIdentifier(action) });
      }
    }
  } else if (idOrAction) {
    crumbs.push({
      label: LABELS[idOrAction] ?? humanizeIdentifier(idOrAction),
      href: action ? `/teams/${teamID}/${section}/${idOrAction}` : undefined,
    });

    if (action) {
      crumbs.push({ label: LABELS[action] ?? humanizeIdentifier(action) });
    }
  }

  return (
    <>
      {crumbs.map((crumb, index) => {
        const isLast = index === crumbs.length - 1;

        return (
          <Fragment key={`${crumb.label}-${index}`}>
            {index > 0 ? (
              <BreadcrumbSeparator className="hidden md:block" />
            ) : null}
            <BreadcrumbItem>
              {isLast || !crumb.href ? (
                <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
              ) : (
                <BreadcrumbLink href={crumb.href}>{crumb.label}</BreadcrumbLink>
              )}
            </BreadcrumbItem>
          </Fragment>
        );
      })}
    </>
  );
}
