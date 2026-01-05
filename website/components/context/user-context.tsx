"use client";

import * as React from "react";

export type UserSummary = {
  id: string;
  displayName: string;
  email: string;
  avatar?: string | null;
};

const UserContext = React.createContext<UserSummary | null>(null);

export function UserProvider({
  user,
  children,
}: {
  user: UserSummary;
  children: React.ReactNode;
}) {
  return <UserContext.Provider value={user}>{children}</UserContext.Provider>;
}

export function useUser() {
  const context = React.useContext(UserContext);
  if (!context) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}
