import * as React from "react"

import { cn } from "@/lib/utils"

export type TimelineItem = {
  id?: React.Key
  label?: React.ReactNode
  meta?: React.ReactNode
  title?: React.ReactNode
  description?: React.ReactNode
  children?: React.ReactNode
}

type TimelineProps = {
  items: TimelineItem[]
  reverse?: boolean
  className?: string
  compact?: boolean
}

export function Timeline({
  items,
  reverse = false,
  className,
  compact = false,
}: TimelineProps) {
  const orderedItems = React.useMemo(
    () => (reverse ? [...items].reverse() : items),
    [items, reverse]
  )

  return (
    <div className={cn("relative", className)}>
      {orderedItems.map((item, index) => (
        <div key={item.id ?? index} className="group relative">
          <div className="flex items-start">
            <div
              className={cn(
                "relative border-l-2 pl-6 sm:pl-8 w-full",
                compact ? "pb-6 group-last:pb-3 space-y-1.5" : "pb-10 group-last:pb-4 space-y-2"
              )}
            >
              <div
                className={cn(
                  "absolute -translate-x-1/2 -left-px rounded-full border-2 border-primary bg-background",
                  compact ? "h-2.5 w-2.5 top-3" : "h-3 w-3 top-4"
                )}
              />

              <div className="space-y-2 min-w-0">
                {item.title ? (
                  <div
                    className={cn(
                      "font-semibold tracking-[-0.01em]",
                      compact ? "text-base" : "text-lg"
                    )}
                  >
                    {item.title}
                  </div>
                ) : null}

                {(item.meta || item.label) && (
                  <div
                    className={cn(
                      "flex flex-wrap items-center gap-2",
                      compact ? "text-xs" : "text-sm"
                    )}
                  >
                    {item.meta ? (
                      <div className="leading-tight">{item.meta}</div>
                    ) : null}
                    {item.label ? (
                      <div
                        className={cn(
                          "text-muted-foreground",
                          compact ? "text-[11px]" : "text-xs sm:text-sm"
                        )}
                      >
                        {item.label}
                      </div>
                    ) : null}
                  </div>
                )}

                {item.description ? (
                  <div
                    className={cn(
                      "text-muted-foreground whitespace-pre-wrap break-words",
                      compact ? "text-sm" : "text-sm sm:text-base"
                    )}
                  >
                    {item.description}
                  </div>
                ) : null}

                {item.children}
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
