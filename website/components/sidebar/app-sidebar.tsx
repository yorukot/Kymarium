"use client";

import * as React from "react";
import {
  Activity,
  LayoutTemplate,
  Settings,
  TriangleAlert,
  Users,
} from "lucide-react";

import { NavMain } from "@/components/sidebar/nav-main";
import { NavUser } from "@/components/sidebar/nav-user";
import { TeamSwitcher } from "@/components/sidebar/team-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useParams, usePathname } from "next/navigation";

// This is sample data.
const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Monitors",
      url: "/teams/{teamID}/monitors",
      icon: Activity,
      isActive: true,
    },
    {
      title: "Incident",
      url: "/teams/{teamID}/incidents",
      icon: TriangleAlert,
    },
    {
      title: "Status Pages",
      url: "/teams/{teamID}/status-pages",
      icon: LayoutTemplate,
    },
    {
      title: "Members",
      url: "/teams/{teamID}/members",
      icon: Users,
    },
    {
      title: "Setting",
      url: "/teams/{teamID}/settings",
      icon: Settings,
    },
  ],
};

function isActivePath(pathname: string, href: string) {
  const p = pathname.replace(/\/$/, "");
  const h = href.replace(/\/$/, "");

  return p === h || p.startsWith(h + "/");
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const params = useParams<{ teamID: string }>();
  const pathname = usePathname()

  const navItems = React.useMemo(() => {
    const teamID = params.teamID;
    const items = teamID
      ? data.navMain.map((item) => ({
          ...item,
          url: item.url.replace("{teamID}", teamID),
        }))
      : data.navMain;

    return items.map((item) => ({
      ...item,
      isActive: isActivePath(pathname, item.url),
    }));
  }, [params.teamID, pathname]);

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navItems} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
