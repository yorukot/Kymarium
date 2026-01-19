import { Timeline as TimelineList, TimelineItem } from "./ui/timeline"

const changelog: TimelineItem[] = [
  {
    title: "Initial Release",
    description:
      "Launched the first version with core features and basic UI components.",
    label: "2025-03-01",
  },
  {
    title: "UI Enhancements",
    description:
      "Improved the user interface with better accessibility and design consistency.",
    label: "2025-03-05",
  },
  {
    title: "Performance Optimization",
    description:
      "Reduced load times and improved overall application performance.",
    label: "2025-03-10",
  },
  {
    title: "New Feature: Dark Mode",
    description:
      "Added support for dark mode, allowing users to switch themes seamlessly.",
    label: "2025-03-15",
  },
  {
    title: "Bug Fixes & Security Patch",
    description:
      "Fixed various minor bugs and patched security vulnerabilities.",
    label: "2025-03-18",
  },
  {
    title: "New Components Added",
    description:
      "Introduced new UI components for better customization and flexibility.",
    label: "2025-03-22",
  },
  {
    title: "Major Update: API Integration",
    description:
      "Integrated external APIs to enhance functionality and data synchronization.",
    label: "2025-04-01",
  },
];

export default function Timeline() {
  return (
    <div className="max-w-(--breakpoint-sm) md:mx-auto py-12 md:py-20 px-6">
      <TimelineList items={changelog} reverse />
    </div>
  );
}
