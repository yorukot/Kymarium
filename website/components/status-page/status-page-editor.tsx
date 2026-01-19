"use client";

import {
  ChartLine,
  CircleDot,
  GripVertical,
  Plus,
  Save,
  Trash2,
} from "lucide-react";
import {
  closestCenter,
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  DragOverEvent,
  KeyboardSensor,
  PointerSensor,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { cn } from "@/lib/utils";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { ApiError } from "@/lib/api/client";
import { updateStatusPage } from "@/lib/api/status-page";
import type { MonitorListItem } from "@/lib/schemas/monitor";
import type {
  StatusPageDetailItem,
  StatusPageElementRawData,
  StatusPageElementType,
} from "@/lib/schemas/status-page";

const ROOT_CONTAINER_ID = "container:root";
const containerID = (key: string) => `container:${key}`;

const statusPageBasicsSchema = z.object({
  title: z.string().min(1, "Title is required.").max(255, "Title is too long."),
  slug: z
    .string()
    .min(3, "Slug must be at least 3 characters.")
    .max(255, "Slug is too long.")
    .regex(/^[a-zA-Z0-9]+$/, "Slug must be alphanumeric (no dashes)."),
});

type StatusPageBasicsValues = z.infer<typeof statusPageBasicsSchema>;

type MonitorItem = {
  uid: string;
  monitorId: string;
  name: string;
  type: StatusPageElementType;
};

type EditorElement =
  | {
      uid: string;
      kind: "group";
      name: string;
      type: StatusPageElementType;
      monitors: MonitorItem[];
    }
  | ({
      uid: string;
      kind: "monitor";
    } & MonitorItem);

function validateElementNames(
  elements: EditorElement[],
): Record<string, string> {
  const errors: Record<string, string> = {};

  for (const element of elements) {
    if (element.kind === "group") {
      if (!element.name.trim()) {
        errors[element.uid] = "Group name is required.";
      }
      for (const monitor of element.monitors) {
        if (!monitor.name.trim()) {
          errors[monitor.uid] = "Monitor name is required.";
        }
      }
      continue;
    }

    if (!element.name.trim()) {
      errors[element.uid] = "Monitor name is required.";
    }
  }

  return errors;
}

function uid(prefix: string) {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return `${prefix}_${crypto.randomUUID()}`;
  }
  return `${prefix}_${Math.random().toString(36).slice(2)}`;
}

function toMonitorItem(raw: {
  monitor_id: string;
  name: string;
  type: StatusPageElementType;
}): MonitorItem {
  return {
    uid: uid("monitor"),
    monitorId: raw.monitor_id,
    name: raw.name,
    type: raw.type,
  };
}

function fromStatusPageElements(
  elements: StatusPageElementRawData[],
): EditorElement[] {
  return (elements ?? []).map((element) => {
    if (element.monitor && element.monitor_id) {
      return {
        uid: uid("element"),
        kind: "monitor",
        monitorId: element.monitor_id,
        name: element.name,
        type: element.type,
      };
    }

    const monitors = (element.monitors ?? []).map((m) =>
      toMonitorItem({ monitor_id: m.monitor_id, name: m.name, type: m.type }),
    );

    return {
      uid: uid("element"),
      kind: "group",
      name: element.name,
      type: element.type,
      monitors,
    };
  });
}

function monitorItemToElement(
  item: MonitorItem,
): Extract<EditorElement, { kind: "monitor" }> {
  return {
    uid: item.uid,
    kind: "monitor",
    monitorId: item.monitorId,
    name: item.name,
    type: item.type,
  };
}

function isContainerID(id: unknown): id is string {
  return typeof id === "string" && id.startsWith("container:");
}

function containerKeyFromID(id: string): string | null {
  return isContainerID(id) ? id.slice("container:".length) : null;
}

function findContainerKey(
  elements: EditorElement[],
  id: string,
): string | null {
  const containerKey = containerKeyFromID(id);
  if (containerKey) return containerKey;

  if (elements.some((el) => el.uid === id)) return "root";

  for (const element of elements) {
    if (element.kind !== "group") continue;
    if (element.monitors.some((m) => m.uid === id)) return element.uid;
  }

  return null;
}

function findItemKind(
  elements: EditorElement[],
  id: string,
): "group" | "monitor" | null {
  for (const element of elements) {
    if (element.uid === id) return element.kind;
    if (element.kind === "group" && element.monitors.some((m) => m.uid === id))
      return "monitor";
  }
  return null;
}

function getOverContainerKey(
  elements: EditorElement[],
  overID: unknown,
): string | null {
  if (!overID) return null;
  return findContainerKey(elements, String(overID));
}

function rootIndex(elements: EditorElement[], id: string): number {
  return elements.findIndex((el) => el.uid === id);
}

function groupIndex(elements: EditorElement[], groupUID: string): number {
  return elements.findIndex((el) => el.kind === "group" && el.uid === groupUID);
}

function groupMonitorIndex(
  elements: EditorElement[],
  groupUID: string,
  monitorUID: string,
): number {
  const idx = groupIndex(elements, groupUID);
  if (idx === -1) return -1;
  const group = elements[idx] as Extract<EditorElement, { kind: "group" }>;
  return group.monitors.findIndex((m) => m.uid === monitorUID);
}

function getInsertionIndex(
  elements: EditorElement[],
  containerKey: string,
  overID: unknown,
): number {
  if (!overID || isContainerID(overID)) {
    if (containerKey === "root") return elements.length;
    const idx = groupIndex(elements, containerKey);
    if (idx === -1) return 0;
    const group = elements[idx] as Extract<EditorElement, { kind: "group" }>;
    return group.monitors.length;
  }

  const over = String(overID);
  if (containerKey === "root") {
    const idx = rootIndex(elements, over);
    return idx === -1 ? elements.length : idx;
  }

  const idx = groupMonitorIndex(elements, containerKey, over);
  if (idx === -1) {
    const groupIdx = groupIndex(elements, containerKey);
    if (groupIdx === -1) return 0;
    const group = elements[groupIdx] as Extract<
      EditorElement,
      { kind: "group" }
    >;
    return group.monitors.length;
  }
  return idx;
}

function moveMonitorBetweenContainers(
  elements: EditorElement[],
  activeID: string,
  fromKey: string,
  toKey: string,
  overID: unknown,
): EditorElement[] {
  if (fromKey === toKey) return elements;

  let monitor: MonitorItem | null = null;
  const next = [...elements];

  if (fromKey === "root") {
    const idx = rootIndex(next, activeID);
    if (idx === -1) return elements;
    const el = next[idx];
    if (!el || el.kind !== "monitor") return elements;
    monitor = {
      uid: el.uid,
      monitorId: el.monitorId,
      name: el.name,
      type: el.type,
    };
    next.splice(idx, 1);
  } else {
    const gIdx = groupIndex(next, fromKey);
    if (gIdx === -1) return elements;
    const group = next[gIdx] as Extract<EditorElement, { kind: "group" }>;
    const mIdx = group.monitors.findIndex((m) => m.uid === activeID);
    if (mIdx === -1) return elements;
    monitor = group.monitors[mIdx];
    const monitors = [...group.monitors];
    monitors.splice(mIdx, 1);
    next[gIdx] = { ...group, monitors };
  }

  if (!monitor) return elements;

  const insertAt = getInsertionIndex(next, toKey, overID);
  if (toKey === "root") {
    const item = monitorItemToElement(monitor);
    next.splice(Math.max(0, Math.min(insertAt, next.length)), 0, item);
    return next;
  }

  const gIdx = groupIndex(next, toKey);
  if (gIdx === -1) return elements;
  const group = next[gIdx] as Extract<EditorElement, { kind: "group" }>;
  const monitors = [...group.monitors];
  monitors.splice(Math.max(0, Math.min(insertAt, monitors.length)), 0, monitor);
  next[gIdx] = { ...group, monitors };
  return next;
}

function reorderWithinRoot(
  elements: EditorElement[],
  activeID: string,
  overID: unknown,
): EditorElement[] {
  const oldIndex = rootIndex(elements, activeID);
  if (oldIndex === -1) return elements;

  const newIndex = getInsertionIndex(elements, "root", overID);
  if (newIndex === oldIndex) return elements;

  const bounded = Math.max(0, Math.min(newIndex, elements.length - 1));
  return arrayMove(elements, oldIndex, bounded);
}

function reorderWithinGroup(
  elements: EditorElement[],
  groupUID: string,
  activeID: string,
  overID: unknown,
): EditorElement[] {
  const gIdx = groupIndex(elements, groupUID);
  if (gIdx === -1) return elements;
  const group = elements[gIdx] as Extract<EditorElement, { kind: "group" }>;

  const oldIndex = group.monitors.findIndex((m) => m.uid === activeID);
  if (oldIndex === -1) return elements;

  const newIndex = getInsertionIndex(elements, groupUID, overID);
  if (newIndex === oldIndex) return elements;

  const bounded = Math.max(0, Math.min(newIndex, group.monitors.length - 1));
  const monitors = arrayMove(group.monitors, oldIndex, bounded);
  const next = [...elements];
  next[gIdx] = { ...group, monitors };
  return next;
}

function buildUpsertElements(elements: EditorElement[]) {
  return elements.map((element, index) => {
    const sortOrder = index + 1;
    if (element.kind === "monitor") {
      return {
        name: element.name,
        type: element.type,
        sortOrder,
        monitor: true,
        monitorId: element.monitorId,
        monitors: [],
      };
    }

    return {
      name: element.name,
      type: element.type,
      sortOrder,
      monitor: false,
      monitors: element.monitors.map((m, mIndex) => ({
        monitorId: m.monitorId,
        name: m.name,
        type: m.type,
        sortOrder: mIndex + 1,
      })),
    };
  });
}

export default function StatusPageEditor({
  teamID,
  statusPage,
  monitorOptions,
}: {
  teamID: string;
  statusPage: StatusPageDetailItem;
  monitorOptions: MonitorListItem[];
}) {
  const [elements, setElements] = useState<EditorElement[]>(() =>
    fromStatusPageElements(statusPage.elements),
  );
  const [activeID, setActiveID] = useState<string | null>(null);
  const [overContainerKey, setOverContainerKey] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  const monitorByID = useMemo(() => {
    return new Map(monitorOptions.map((monitor) => [monitor.id, monitor]));
  }, [monitorOptions]);

  const basicsForm = useForm<StatusPageBasicsValues>({
    resolver: zodResolver(statusPageBasicsSchema),
    defaultValues: {
      title: statusPage.title,
      slug: statusPage.slug,
    },
    mode: "onSubmit",
  });

  const title = basicsForm.watch("title");
  const slug = basicsForm.watch("slug");

  const [nameErrors, setNameErrors] = useState<Record<string, string>>({});

  const setNameError = (uid: string, message: string | null) => {
    setNameErrors((prev) => {
      if (!message) {
        if (!(uid in prev)) return prev;
        const next = { ...prev };
        delete next[uid];
        return next;
      }

      if (prev[uid] === message) return prev;
      return { ...prev, [uid]: message };
    });
  };

  const [pickerOpen, setPickerOpen] = useState(false);
  const [pickerTarget, setPickerTarget] = useState<DropTarget | null>(null);

  const [confirmDeleteGroupOpen, setConfirmDeleteGroupOpen] = useState(false);
  const [confirmDeleteGroupUid, setConfirmDeleteGroupUid] = useState<
    string | null
  >(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  const rootItems = useMemo(() => elements.map((e) => e.uid), [elements]);
  const { setNodeRef: setRootNodeRef } = useDroppable({
    id: ROOT_CONTAINER_ID,
  });

  const openMonitorPicker = (target: DropTarget) => {
    setPickerTarget(target);
    setPickerOpen(true);
  };

  type DropTarget =
    | { container: "root"; index: number }
    | { container: "group"; groupUid: string; index: number };

  const addGroup = () => {
    setElements((prev) => [
      ...prev,
      {
        uid: uid("element"),
        kind: "group",
        name: "New group",
        type: "current_status_indicator",
        monitors: [],
      },
    ]);
  };

  const addMonitor = (monitor: MonitorListItem, target: DropTarget) => {
    const item: MonitorItem = {
      uid: uid("element"),
      monitorId: monitor.id,
      name: monitor.name,
      type: "current_status_indicator",
    };

    setElements((prev) => {
      if (target.container === "root") {
        const next = [...prev];
        next.splice(
          Math.max(0, Math.min(target.index, next.length)),
          0,
          monitorItemToElement(item),
        );
        return next;
      }

      const groupIndex = prev.findIndex(
        (e) => e.kind === "group" && e.uid === target.groupUid,
      );
      if (groupIndex === -1) return prev;
      const group = prev[groupIndex] as Extract<
        EditorElement,
        { kind: "group" }
      >;
      const monitors = [...group.monitors];
      monitors.splice(
        Math.max(0, Math.min(target.index, monitors.length)),
        0,
        item,
      );
      const next = [...prev];
      next[groupIndex] = { ...group, monitors };
      return next;
    });
  };

  const requestDeleteGroup = (groupUid: string) => {
    setConfirmDeleteGroupUid(groupUid);
    setConfirmDeleteGroupOpen(true);
  };

  const deleteGroup = () => {
    const groupUid = confirmDeleteGroupUid;
    setConfirmDeleteGroupOpen(false);
    setConfirmDeleteGroupUid(null);
    if (!groupUid) return;

    setNameError(groupUid, null);

    setElements((prev) => {
      const idx = prev.findIndex(
        (e) => e.kind === "group" && e.uid === groupUid,
      );
      if (idx === -1) return prev;
      const group = prev[idx] as Extract<EditorElement, { kind: "group" }>;
      const ungrouped = group.monitors.map(monitorItemToElement);
      const next = [...prev];
      next.splice(idx, 1, ...ungrouped);
      return next;
    });
  };

  const deleteMonitorElement = (monitorUid: string) => {
    setNameError(monitorUid, null);
    setElements((prev) =>
      prev.filter((e) => !(e.kind === "monitor" && e.uid === monitorUid)),
    );
  };

  const deleteGroupMonitor = (groupUid: string, monitorUid: string) => {
    setNameError(monitorUid, null);
    setElements((prev) => {
      const idx = prev.findIndex(
        (e) => e.kind === "group" && e.uid === groupUid,
      );
      if (idx === -1) return prev;
      const group = prev[idx] as Extract<EditorElement, { kind: "group" }>;
      const nextMonitors = group.monitors.filter((m) => m.uid !== monitorUid);
      const next = [...prev];
      next[idx] = { ...group, monitors: nextMonitors };
      return next;
    });
  };

  const save = basicsForm.handleSubmit(
    async (values) => {
      const nextNameErrors = validateElementNames(elements);
      if (Object.keys(nextNameErrors).length > 0) {
        setNameErrors(nextNameErrors);
        toast.error("Fix the element name errors before saving.");
        return;
      }

      setIsSaving(true);
      try {
        const res = await updateStatusPage(teamID, statusPage.id, {
          title: values.title.trim(),
          slug: values.slug.trim(),
          elements: buildUpsertElements(elements),
        });
        if (res.data?.data) {
          basicsForm.reset({
            title: res.data.data.status_page.title,
            slug: res.data.data.status_page.slug,
          });
          setElements(fromStatusPageElements(res.data.data.elements));
        }
        toast.success("Status page saved.");
      } catch (error) {
        if (error instanceof ApiError) {
          toast.error(
            error.status >= 500
              ? "Server error. Please try again later."
              : error.message,
          );
        } else {
          toast.error("Network error. Please try again.");
        }
      } finally {
        setIsSaving(false);
      }
    },
    () => {
      toast.error("Fix the validation errors before saving.");
    },
  );

  const apiBase = process.env.NEXT_PUBLIC_FRONTEND_URL;
  const slugForUrl = slug?.trim() || statusPage.slug;
  const publicUrl = apiBase ? `${apiBase}/s/${slugForUrl}` : null;
  const activeOverlay = useMemo(() => {
    if (!activeID) return null;
    for (const el of elements) {
      if (el.uid === activeID) return el;
      if (el.kind === "group") {
        const m = el.monitors.find((x) => x.uid === activeID);
        if (m) return { ...m, kind: "monitor" as const };
      }
    }
    return null;
  }, [activeID, elements]);

  const handleDragStart = (event: DragStartEvent) => {
    const id = String(event.active.id);
    setActiveID(id);
  };

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over } = event;
    if (!over) {
      setOverContainerKey(null);
      return;
    }

    const activeId = String(active.id);
    const activeKind = findItemKind(elements, activeId);
    if (activeKind !== "monitor") {
      setOverContainerKey(null);
      return;
    }

    const overId = String(over.id);
    const overKind = findItemKind(elements, overId);
    const nextOverContainer =
      overKind === "group" ? overId : getOverContainerKey(elements, over.id);

    setOverContainerKey(nextOverContainer);

    // Make cross-container moves easier by moving the item as soon as the pointer
    // crosses into the other container (root <-> group).
    setElements((prev) => {
      const fromKey = findContainerKey(prev, activeId);
      const toKey =
        overKind === "group"
          ? overId
          : (getOverContainerKey(prev, over.id) ?? nextOverContainer);

      if (!fromKey || !toKey) return prev;
      if (fromKey === toKey) return prev;

      return moveMonitorBetweenContainers(
        prev,
        activeId,
        fromKey,
        toKey,
        over.id,
      );
    });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over) {
      setActiveID(null);
      setOverContainerKey(null);
      return;
    }

    const activeId = String(active.id);
    const overId = over.id;

    setElements((prev) => {
      const activeKind = findItemKind(prev, activeId);
      if (!activeKind) return prev;

      const activeContainer = findContainerKey(prev, activeId);
      const overIDString = String(overId);
      const overKind = findItemKind(prev, overIDString);
      const overContainer =
        activeKind === "monitor" && overKind === "group"
          ? overIDString
          : getOverContainerKey(prev, overId);
      if (!activeContainer || !overContainer) return prev;

      // Groups can only be sorted on the root list.
      if (activeKind === "group") {
        const effectiveOverID =
          overContainer === "root" ? overId : overContainer;
        return reorderWithinRoot(prev, activeId, effectiveOverID);
      }

      // Monitors can be moved between root and groups.
      if (activeContainer !== overContainer) {
        const effectiveOverID =
          typeof overId === "string" && findItemKind(prev, overId) === "group"
            ? containerID(overId)
            : overId;

        return moveMonitorBetweenContainers(
          prev,
          activeId,
          activeContainer,
          overContainer,
          effectiveOverID,
        );
      }

      if (overContainer === "root") {
        return reorderWithinRoot(prev, activeId, overId);
      }

      return reorderWithinGroup(prev, overContainer, activeId, overId);
    });

    setActiveID(null);
    setOverContainerKey(null);
  };

  return (
    <div className="flex flex-col gap-4 pb-12">
      <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
        <div className="flex flex-col gap-1 min-w-0">
          <h1 className="text-xl font-semibold truncate">{title}</h1>
          <p className="text-sm text-muted-foreground truncate">/{slug}</p>
        </div>

        <div className="flex items-center gap-2 ">
          <Button onClick={() => void save()} disabled={isSaving}>
            {isSaving ? (
              <>
                <Spinner /> Saving…
              </>
            ) : (
              <>
                <Save /> Save
              </>
            )}
          </Button>
        </div>
      </div>

      {publicUrl ? (
        <div className="rounded-md border px-3 py-2 text-sm">
          Public URL:{" "}
          <a
            href={publicUrl}
            target="_blank"
            rel="noreferrer"
            className="underline underline-offset-4"
          >
            {publicUrl}
          </a>
        </div>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle>Basics</CardTitle>
          <CardDescription>
            Update the public title and slug (slug must be alphanumeric).
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
            <div className="flex flex-col gap-1.5">
              <div className="text-sm font-medium">Title</div>
              <Input
                {...basicsForm.register("title")}
                placeholder="Kymarium Status"
                aria-label="Status page title"
                aria-invalid={!!basicsForm.formState.errors.title}
              />
              {basicsForm.formState.errors.title?.message ? (
                <div className="text-sm text-destructive">
                  {basicsForm.formState.errors.title.message}
                </div>
              ) : null}
            </div>

            <div className="flex flex-col gap-1.5">
              <div className="text-sm font-medium">Slug</div>
              <Input
                {...basicsForm.register("slug")}
                placeholder="kymariumstatus"
                aria-label="Status page slug"
                aria-invalid={!!basicsForm.formState.errors.slug}
              />
              <div className="text-xs text-muted-foreground">
                No dashes; only a-z, A-Z, 0-9.
              </div>
              {basicsForm.formState.errors.slug?.message ? (
                <div className="text-sm text-destructive">
                  {basicsForm.formState.errors.slug.message}
                </div>
              ) : null}
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="flex flex-col gap-2">
        <div className="flex flex-col gap-0.5">
          <h2 className="text-sm font-medium">Layout</h2>
          <p className="text-sm text-muted-foreground">
            Top-level order must be consecutive. Group contents are ordered
            independently.
          </p>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row">
          <Button
            variant="outline"
            onClick={() => addGroup()}
            className="w-full sm:flex-1"
          >
            <Plus />
            Add group
          </Button>
          <Button
            variant="outline"
            className="w-full sm:flex-1"
            onClick={() =>
              openMonitorPicker({ container: "root", index: elements.length })
            }
          >
            <Plus />
            Add monitor
          </Button>
        </div>
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={handleDragStart}
          onDragOver={handleDragOver}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={rootItems}
            strategy={verticalListSortingStrategy}
          >
            <div ref={setRootNodeRef} className="flex flex-col gap-2">
              {elements.length === 0 ? (
                <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
                  No elements yet. Add a group or a monitor.
                </div>
              ) : null}

              {elements.map((element) => {
                if (element.kind === "monitor") {
                  return (
                    <SortableRootMonitor
                      key={element.uid}
                      element={element}
                      originalMonitor={monitorByID.get(element.monitorId) ?? null}
                      onChangeName={(nextName) => {
                        setElements((prev) =>
                          prev.map((e) =>
                            e.uid === element.uid && e.kind === "monitor"
                              ? { ...e, name: nextName }
                              : e,
                          ),
                        );
                        if (nextName.trim()) setNameError(element.uid, null);
                      }}
                      onBlurName={() =>
                        setNameError(
                          element.uid,
                          element.name.trim()
                            ? null
                            : "Monitor name is required.",
                        )
                      }
                      onChangeType={(nextType) =>
                        setElements((prev) =>
                          prev.map((e) =>
                            e.uid === element.uid && e.kind === "monitor"
                              ? { ...e, type: nextType }
                              : e,
                          ),
                        )
                      }
                      onDelete={() => deleteMonitorElement(element.uid)}
                      nameError={nameErrors[element.uid]}
                    />
                  );
                }

                return (
                  <SortableGroup
                    key={element.uid}
                    group={element}
                    overContainerKey={overContainerKey}
                    onChangeName={(nextName) => {
                      setElements((prev) =>
                        prev.map((e) =>
                          e.uid === element.uid && e.kind === "group"
                            ? { ...e, name: nextName }
                            : e,
                        ),
                      );
                      if (nextName.trim()) setNameError(element.uid, null);
                    }}
                    onBlurName={() =>
                      setNameError(
                        element.uid,
                        element.name.trim() ? null : "Group name is required.",
                      )
                    }
                    onChangeType={(nextType) =>
                      setElements((prev) =>
                        prev.map((e) =>
                          e.uid === element.uid && e.kind === "group"
                            ? { ...e, type: nextType }
                            : e,
                        ),
                      )
                    }
                    onAddMonitor={() =>
                      openMonitorPicker({
                        container: "group",
                        groupUid: element.uid,
                        index: element.monitors.length,
                      })
                    }
                    onDelete={() => requestDeleteGroup(element.uid)}
                    onDeleteMonitor={(monitorUid) =>
                      deleteGroupMonitor(element.uid, monitorUid)
                    }
                    onChangeMonitorName={(monitorUid, nextName) => {
                      setElements((prev) =>
                        prev.map((e) => {
                          if (e.kind !== "group" || e.uid !== element.uid)
                            return e;
                          return {
                            ...e,
                            monitors: e.monitors.map((m) =>
                              m.uid === monitorUid
                                ? { ...m, name: nextName }
                                : m,
                            ),
                          };
                        }),
                      );
                      if (nextName.trim()) setNameError(monitorUid, null);
                    }}
                    onBlurMonitorName={(monitorUid, currentName) =>
                      setNameError(
                        monitorUid,
                        currentName.trim() ? null : "Monitor name is required.",
                      )
                    }
                    onChangeMonitorType={(monitorUid, nextType) =>
                      setElements((prev) =>
                        prev.map((e) => {
                          if (e.kind !== "group" || e.uid !== element.uid)
                            return e;
                          return {
                            ...e,
                            monitors: e.monitors.map((m) =>
                              m.uid === monitorUid
                                ? { ...m, type: nextType }
                                : m,
                            ),
                          };
                        }),
                      )
                    }
                    groupNameError={nameErrors[element.uid]}
                    monitorNameErrors={nameErrors}
                    monitorByID={monitorByID}
                  />
                );
              })}
            </div>
          </SortableContext>

          <DragOverlay>
            {activeOverlay ? (
              <OverlayItem
                label={
                  activeOverlay.kind === "group"
                    ? activeOverlay.name
                    : activeOverlay.name
                }
              />
            ) : null}
          </DragOverlay>
        </DndContext>
      </div>

      <MonitorPickerDialog
        open={pickerOpen}
        onOpenChange={setPickerOpen}
        title="Add a monitor"
        description="Search and pick a monitor to add to this status page."
        options={monitorOptions}
        onPick={(m) => {
          if (!pickerTarget) return;
          addMonitor(m, pickerTarget);
          setPickerOpen(false);
          setPickerTarget(null);
        }}
      />

      <Dialog
        open={confirmDeleteGroupOpen}
        onOpenChange={setConfirmDeleteGroupOpen}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete group?</DialogTitle>
            <DialogDescription>
              This will remove the group container and ungroup its monitors
              (they will become top-level items).
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setConfirmDeleteGroupOpen(false)}
            >
              Cancel
            </Button>
            <Button variant="destructive" onClick={() => deleteGroup()}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

function MonitorPickerDialog({
  open,
  onOpenChange,
  title,
  description,
  options,
  onPick,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  options: MonitorListItem[];
  onPick: (monitor: MonitorListItem) => void;
}) {
  const sorted = useMemo(() => {
    return [...options].sort((a, b) => a.name.localeCompare(b.name));
  }, [options]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="p-0 overflow-hidden">
        <DialogHeader className="px-6 pt-6">
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>

        <div className="px-4 pb-4">
          <Command>
            <CommandInput placeholder="Search monitors..." />
            <CommandList>
              <CommandEmpty>No monitors found.</CommandEmpty>
              <CommandGroup heading="Monitors">
                {sorted.map((monitor) => {
                  return (
                    <CommandItem
                      key={monitor.id}
                      value={monitor.name}
                      onSelect={() => {
                        onPick(monitor);
                      }}
                    >
                      <span className="min-w-0 flex-1 truncate">
                        {monitor.name}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {monitor.type}
                      </span>
                    </CommandItem>
                  );
                })}
              </CommandGroup>
            </CommandList>
          </Command>
        </div>
      </DialogContent>
    </Dialog>
  );
}

function OverlayItem({ label }: { label: string }) {
  return (
    <div className="rounded-md border bg-card px-3 py-2 shadow-sm text-sm">
      {label}
    </div>
  );
}

function ElementTypeSelect({
  value,
  onChange,
  ariaLabel,
}: {
  value: StatusPageElementType;
  onChange: (value: StatusPageElementType) => void;
  ariaLabel: string;
}) {
  const options: Array<{
    value: StatusPageElementType;
    label: string;
    Icon: typeof CircleDot;
  }> = [
    {
      value: "current_status_indicator",
      label: "Current status indicator",
      Icon: CircleDot,
    },
    {
      value: "historical_timeline",
      label: "Historical timeline",
      Icon: ChartLine,
    },
  ];

  return (
    <Select
      value={value}
      onValueChange={(v) => onChange(v as StatusPageElementType)}
    >
      <SelectTrigger
        size="sm"
        className="w-14 px-2 sm:w-55 sm:px-3"
        aria-label={ariaLabel}
        onPointerDown={(event) => event.stopPropagation()}
        onKeyDownCapture={(event) => event.stopPropagation()}
      >
        <SelectValue className="sr-only" />
      </SelectTrigger>
      <SelectContent>
        {options.map(({ value: optionValue, label, Icon }) => (
          <SelectItem key={optionValue} value={optionValue}>
            <span className="flex items-center gap-2">
              <Icon className="size-4" />
              {label}
            </span>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

function SortableHandle({
  setActivatorNodeRef,
  listeners,
  attributes,
}: {
  setActivatorNodeRef: (element: HTMLElement | null) => void;
  listeners: ReturnType<typeof useSortable>["listeners"];
  attributes: ReturnType<typeof useSortable>["attributes"];
}) {
  return (
    <button
      type="button"
      ref={setActivatorNodeRef}
      {...attributes}
      {...listeners}
      className="text-muted-foreground cursor-grab active:cursor-grabbing"
      aria-label="Drag handle"
    >
      <GripVertical size={16} />
    </button>
  );
}

function SortableRootMonitor({
  element,
  originalMonitor,
  onChangeName,
  onBlurName,
  onChangeType,
  onDelete,
  nameError,
}: {
  element: Extract<EditorElement, { kind: "monitor" }>;
  originalMonitor: MonitorListItem | null;
  onChangeName: (name: string) => void;
  onBlurName: () => void;
  onChangeType: (type: StatusPageElementType) => void;
  onDelete: () => void;
  nameError?: string;
}) {
  const {
    setNodeRef,
    attributes,
    listeners,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: element.uid,
    data: { kind: "monitor", container: "root" },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "rounded-md border bg-card px-3 py-2",
        isDragging && "opacity-60",
      )}
    >
      <div className="flex items-center gap-2">
        <SortableHandle
          setActivatorNodeRef={setActivatorNodeRef}
          attributes={attributes}
          listeners={listeners}
        />
        <Input
          value={element.name}
          onChange={(e) => onChangeName(e.target.value)}
          className="flex-1"
          aria-label="Monitor name"
          onBlur={onBlurName}
          aria-invalid={!!nameError}
        />
        <ElementTypeSelect
          value={element.type}
          onChange={onChangeType}
          ariaLabel="Monitor element type"
        />
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Remove monitor"
          onClick={onDelete}
        >
          <Trash2 size={16} />
        </Button>
      </div>
      {originalMonitor ? (
        <div
          className="mt-1 min-w-0 truncate text-xs text-muted-foreground"
          title={`${originalMonitor.name}: ${originalMonitor.targetValue}`}
        >
          {originalMonitor.name}: {originalMonitor.targetValue}
        </div>
      ) : null}
      {nameError ? (
        <div className="mt-1 text-sm text-destructive">{nameError}</div>
      ) : null}
    </div>
  );
}

function SortableGroup({
  group,
  overContainerKey,
  onChangeName,
  onBlurName,
  onChangeType,
  onAddMonitor,
  onDelete,
  onDeleteMonitor,
  onChangeMonitorName,
  onBlurMonitorName,
  onChangeMonitorType,
  groupNameError,
  monitorNameErrors,
  monitorByID,
}: {
  group: Extract<EditorElement, { kind: "group" }>;
  overContainerKey: string | null;
  onChangeName: (name: string) => void;
  onBlurName: () => void;
  onChangeType: (type: StatusPageElementType) => void;
  onAddMonitor: () => void;
  onDelete: () => void;
  onDeleteMonitor: (monitorUid: string) => void;
  onChangeMonitorName: (monitorUid: string, name: string) => void;
  onBlurMonitorName: (monitorUid: string, currentName: string) => void;
  onChangeMonitorType: (
    monitorUid: string,
    type: StatusPageElementType,
  ) => void;
  groupNameError?: string;
  monitorNameErrors: Record<string, string>;
  monitorByID: Map<string, MonitorListItem>;
}) {
  const {
    setNodeRef,
    attributes,
    listeners,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: group.uid,
    data: { kind: "group", container: "root" },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const groupContainerID = containerID(group.uid);
  const monitorItems = group.monitors.map((m) => m.uid);
  const { setNodeRef: setMonitorsNodeRef, isOver } = useDroppable({
    id: groupContainerID,
  });
  const isHighlighted = isOver || overContainerKey === group.uid;

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn("rounded-md border bg-card", isDragging && "opacity-60")}
    >
      <div className="flex items-center gap-2 px-3 py-2">
        <SortableHandle
          setActivatorNodeRef={setActivatorNodeRef}
          attributes={attributes}
          listeners={listeners}
        />
        <Input
          value={group.name}
          onChange={(e) => onChangeName(e.target.value)}
          className="flex-1"
          aria-label="Group name"
          onBlur={onBlurName}
          aria-invalid={!!groupNameError}
        />
        <ElementTypeSelect
          value={group.type}
          onChange={onChangeType}
          ariaLabel="Group element type"
        />
        <div className="flex items-center gap-1 shrink-0">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={onAddMonitor}
            aria-label="Add monitor to group"
          >
            <Plus />
            <span className="sr-only sm:not-sr-only sm:inline">Add</span>
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            aria-label="Delete group"
            onClick={onDelete}
          >
            <Trash2 size={16} />
          </Button>
        </div>
      </div>

      {groupNameError ? (
        <div className="px-3 pb-2 text-sm text-destructive">
          {groupNameError}
        </div>
      ) : null}

      <div className="border-t px-3 py-2">
        <div className="text-sm text-muted-foreground mb-2">Monitors</div>

        <SortableContext
          items={monitorItems}
          strategy={verticalListSortingStrategy}
        >
          <div
            ref={setMonitorsNodeRef}
            className={cn(
              "flex flex-col gap-2 rounded-md p-2 transition-colors min-h-14",
              "border border-dashed",
              isHighlighted
                ? "border-primary bg-primary/5"
                : "border-transparent",
            )}
          >
            {group.monitors.length === 0 ? (
              <div className="rounded-md border border-dashed px-3 py-2 text-sm text-muted-foreground">
                Drag a monitor here or click “Add”.
              </div>
            ) : null}

            {group.monitors.map((monitor) => (
              <SortableGroupMonitor
                key={monitor.uid}
                monitor={monitor}
                groupUid={group.uid}
                originalMonitor={monitorByID.get(monitor.monitorId) ?? null}
                onChangeName={(name) => onChangeMonitorName(monitor.uid, name)}
                onBlurName={() => onBlurMonitorName(monitor.uid, monitor.name)}
                onChangeType={(type) => onChangeMonitorType(monitor.uid, type)}
                onDelete={() => onDeleteMonitor(monitor.uid)}
                nameError={monitorNameErrors[monitor.uid]}
              />
            ))}
          </div>
        </SortableContext>
      </div>
    </div>
  );
}

function SortableGroupMonitor({
  monitor,
  groupUid,
  originalMonitor,
  onChangeName,
  onBlurName,
  onChangeType,
  onDelete,
  nameError,
}: {
  monitor: MonitorItem;
  groupUid: string;
  originalMonitor: MonitorListItem | null;
  onChangeName: (name: string) => void;
  onBlurName: () => void;
  onChangeType: (type: StatusPageElementType) => void;
  onDelete: () => void;
  nameError?: string;
}) {
  const {
    setNodeRef,
    attributes,
    listeners,
    setActivatorNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: monitor.uid,
    data: { kind: "monitor", container: groupUid },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "rounded-md border bg-background px-3 py-2",
        isDragging && "opacity-60",
      )}
    >
      <div className="flex items-center gap-2">
        <SortableHandle
          setActivatorNodeRef={setActivatorNodeRef}
          attributes={attributes}
          listeners={listeners}
        />
        <Input
          value={monitor.name}
          onChange={(e) => onChangeName(e.target.value)}
          className="flex-1"
          aria-label="Monitor name"
          onBlur={onBlurName}
          aria-invalid={!!nameError}
        />
        <ElementTypeSelect
          value={monitor.type}
          onChange={onChangeType}
          ariaLabel="Monitor element type"
        />
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Remove monitor"
          onClick={onDelete}
        >
          <Trash2 size={16} />
        </Button>
      </div>
      {originalMonitor ? (
        <div
          className="mt-1 min-w-0 truncate text-xs text-muted-foreground"
          title={`${originalMonitor.name}: ${originalMonitor.targetValue}`}
        >
          {originalMonitor.name}: {originalMonitor.targetValue}
        </div>
      ) : null}
      {nameError ? (
        <div className="mt-1 text-sm text-destructive">{nameError}</div>
      ) : null}
    </div>
  );
}
