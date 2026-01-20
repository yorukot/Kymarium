import type { Metadata } from "next";
import AccountSettingsPageClient from "@/components/user/account-settings-page";

export const metadata: Metadata = {
  title: "Account Settings",
};

export default function AccountPage() {
  return <AccountSettingsPageClient />;
}
