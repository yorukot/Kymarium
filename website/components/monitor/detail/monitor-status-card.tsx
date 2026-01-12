import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";
import { Card, CardContent } from "@/components/ui/card";

const monitorStatusCardVariants = cva("p-3 rounded-md", {
  variants: {
    tone: {
      default: "",
      successed: "bg-successed/20 border-successed",
      destructive: "bg-destructive/20 border-destructive",
    },
  },
  defaultVariants: {
    tone: "default",
  },
});

const monitorStatusCardTitleVariants = cva("text-sm", {
  variants: {
    tone: {
      default: "text-muted-foreground",
      successed: "text-successed",
      destructive: "text-destructive",
    },
  },
  defaultVariants: {
    tone: "default",
  },
});

const monitorStatusCardValueVariants = cva("text-md", {
  variants: {
    tone: {
      default: "text-foreground",
      successed: "text-successed",
      destructive: "text-destructive",
    },
  },
  defaultVariants: {
    tone: "default",
  },
});

export type MonitorStatusCardTone =
  | NonNullable<VariantProps<typeof monitorStatusCardVariants>["tone"]>
  | "default";

export function MonitorStatusCard({
  title,
  value,
  icon,
  tone = "default",
  valueTone = "default",
  className,
  titleClassName,
  iconClassName,
  valueClassName,
  ...props
}: Omit<React.ComponentProps<"div">, "children"> & {
  title: React.ReactNode;
  value: React.ReactNode;
  icon?: React.ReactNode;
  tone?: MonitorStatusCardTone;
  valueTone?: MonitorStatusCardTone;
  titleClassName?: string;
  iconClassName?: string;
  valueClassName?: string;
}) {
  return (
    <Card
      className={cn(monitorStatusCardVariants({ tone }), className)}
      {...props}
    >
      <CardContent className="p-0">
        <div className="flex items-center gap-1 font-bold">
          {icon ? (
            <span
              className={cn(monitorStatusCardTitleVariants({ tone }), iconClassName)}
            >
              {icon}
            </span>
          ) : null}
          <h3
            className={cn(
              monitorStatusCardTitleVariants({ tone }),
              titleClassName,
            )}
          >
            {title}
          </h3>
        </div>
        <span
          className={cn(
            monitorStatusCardValueVariants({ tone: valueTone }),
            valueClassName,
          )}
        >
          {value}
        </span>
      </CardContent>
    </Card>
  );
}
