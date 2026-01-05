import Monitor from "@/components/monitor/monitor";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";

export default function MonitorsPage() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <div className="flex flex-col gap-1">
          <span className="text-xl font-bold">Monitor List</span>
          <span className="text-sm text-muted-foreground">
            All the monitor are down below
          </span>
        </div>
        <Button>
          <Plus />
          Add Monitor
        </Button>
      </div>
      <Monitor />
    </div>
  );
}
